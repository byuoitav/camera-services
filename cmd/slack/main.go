package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/aver"
	"github.com/byuoitav/axis"
	"github.com/byuoitav/camera-services/couch"
	"github.com/byuoitav/central-event-system/hub/base"
	"github.com/byuoitav/central-event-system/messenger"
	"github.com/byuoitav/common/v2/events"
	"github.com/nlopes/slack"
	"github.com/spf13/pflag"
)

type Snapshotter interface {
	Snapshot(context.Context) (image.Image, error)
}

type ConfigService interface {
	CameraPreset(ctx context.Context, camID, presetID string) (string, error)
}

type data struct {
	client    *slack.Client
	lastTaken map[string]time.Time
	sync.Mutex

	slackChannelID string
	averUsername   string
	averPassword   string
	snapshotDelay  time.Duration

	configService ConfigService
}

func main() {
	var (
		slackToken string
		hubAddress string

		dbAddr     string
		dbUsername string
		dbPassword string
		dbInsecure bool
	)

	d := &data{
		lastTaken: make(map[string]time.Time),
	}

	pflag.StringVar(&dbAddr, "db-address", "", "database address")
	pflag.StringVar(&dbUsername, "db-username", "", "database username")
	pflag.StringVar(&dbPassword, "db-password", "", "database password")
	pflag.BoolVar(&dbInsecure, "db-insecure", false, "don't use SSL in database connection")
	pflag.StringVar(&slackToken, "slack-token", "", "slack token")
	pflag.StringVar(&d.slackChannelID, "channel-id", "", "slack channel id")
	pflag.StringVar(&d.averUsername, "aver-username", "", "aver camera username")
	pflag.StringVar(&d.averPassword, "aver-password", "", "aver camera password")
	pflag.DurationVar(&d.snapshotDelay, "snapshot-delay", 5*time.Second, "snapshot delay (1m5s)")
	pflag.StringVar(&hubAddress, "hub-address", "", "event hub address")
	pflag.Parse()

	d.client = slack.New(slackToken)

	// context for setup
	sctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build the config service
	if dbInsecure {
		dbAddr = "http://" + dbAddr
	} else {
		dbAddr = "https://" + dbAddr
	}

	var csOpts []couch.Option
	if dbUsername != "" {
		csOpts = append(csOpts, couch.WithBasicAuth(dbUsername, dbPassword))
	}

	cs, err := couch.New(sctx, dbAddr, csOpts...)
	if err != nil {
		log.Fatalf("unable to create config service: %s", err)
	}

	d.configService = cs

	// get all of the events
	messenger, nerr := messenger.BuildMessenger(hubAddress, base.Messenger, 4096)
	if nerr != nil {
		log.Fatalf("unable to build messenger: %s", nerr.Error())
	}

	messenger.SubscribeToRooms("*")

	for {
		event := messenger.ReceiveEvent()
		if event.Key != "GoToPreset" {
			continue
		}

		if event.TargetDevice.DeviceID == "" {
			log.Printf("[WARN] invalid event: %+v", event)
			continue
		}

		d.Lock()

		if time.Since(d.lastTaken[event.TargetDevice.DeviceID]).Hours() >= 24 {
			d.lastTaken[event.TargetDevice.DeviceID] = time.Now()
			go d.HandleEvent(event)
		}

		d.Unlock()
	}
}

func (d *data) HandleEvent(event events.Event) {
	fail := func() {
		d.Lock()
		d.lastTaken[event.TargetDevice.DeviceID] = time.Time{}
		d.Unlock()
	}

	var cam Snapshotter
	if strings.Contains(event.GeneratingSystem, "axis") {
		cam = &axis.P5414E{
			Address: event.TargetDevice.DeviceID + ".byu.edu",
		}
	} else if strings.Contains(event.GeneratingSystem, "aver") {
		cam = &aver.Pro520{
			Address:  event.TargetDevice.DeviceID + ".byu.edu",
			Username: d.averUsername,
			Password: d.averPassword,
		}
	} else {
		log.Printf("[WARN] unknown generating system %q", event.GeneratingSystem)
		fail()
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var presetID string
	if v, ok := event.Data.(map[string]interface{}); ok {
		presetID, _ = v["preset"].(string)
	}

	preset, _ := d.configService.CameraPreset(ctx, event.TargetDevice.DeviceID, presetID)

	time.Sleep(d.snapshotDelay)

	if err := d.UploadSnapshot(ctx, cam, event.TargetDevice.DeviceID, preset); err != nil {
		log.Printf("[ERROR] unable to upload screenshot for %q: %s", event.TargetDevice.DeviceID, err)
		fail()
		return
	}

	log.Printf("Successfully uploaded screenshot for %q", event.TargetDevice.DeviceID)
}

func (d *data) UploadSnapshot(ctx context.Context, cam Snapshotter, id, preset string) error {
	snap, err := cam.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("unable to take snapshot: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, snap, nil); err != nil {
		return fmt.Errorf("unable to encode snapshot: %w", err)
	}

	now := time.Now()

	_, err = d.client.UploadFileContext(ctx, slack.FileUploadParameters{
		Filetype: "jpg",
		Filename: fmt.Sprintf("%s_%s.jpg", id, now.Format(time.RFC3339)),
		Title:    fmt.Sprintf("%s-%s snapshot @ %s", id, preset, now.Format(time.RFC3339)),
		Reader:   buf,
		Channels: []string{d.slackChannelID},
	})
	if err != nil {
		return fmt.Errorf("unable to post snapshot: %w", err)
	}

	return nil
}
