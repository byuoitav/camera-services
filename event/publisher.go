package event

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
)

type Publisher struct {
	GeneratingSystem string

	Resolver *net.Resolver
	URL      string
}

type event struct {
	GeneratingSystem string      `json:"generating-system"`
	Timestamp        time.Time   `json:"timestamp"`
	Tags             []string    `json:"event-tags"`
	TargetDevice     deviceInfo  `json:"target-device"`
	AffectedRoom     roomInfo    `json:"affected-room"`
	Key              string      `json:"key"`
	Value            string      `json:"value"`
	User             string      `json:"user"`
	Data             interface{} `json:"data,omitempty"`
}

type roomInfo struct {
	BuildingID string `json:"buildingID,omitempty"`
	RoomID     string `json:"roomID,omitempty"`
}

type deviceInfo struct {
	roomInfo
	DeviceID string `json:"deviceID,omitempty"`
}

func (p *Publisher) Publish(ctx context.Context, info cameraservices.RequestInfo) error {
	info.Data["cameraIP"] = info.CameraIP.String()

	event := event{
		GeneratingSystem: p.GeneratingSystem,
		Timestamp:        info.Timestamp,
		User:             info.SourceIP.String(),
		Key:              info.Action,
		Value:            info.Duration.String(),
		Tags: []string{
			"cameraControl",
		},
		Data: info.Data,
	}

	event = p.handleIPs(ctx, info, event)
	return p.publish(ctx, event)
}

func (p *Publisher) Error(ctx context.Context, err cameraservices.RequestError) error {
	err.Data["cameraIP"] = err.CameraIP.String()
	err.Data["duration"] = err.Duration.String()

	event := event{
		GeneratingSystem: p.GeneratingSystem,
		Timestamp:        err.Timestamp,
		User:             err.SourceIP.String(),
		Key:              err.Action,
		Value:            err.Error,
		Tags: []string{
			"cameraControl",
			"error",
		},
		Data: err.Data,
	}

	event = p.handleIPs(ctx, err.RequestInfo, event)
	return p.publish(ctx, event)
}

func (p *Publisher) handleIPs(ctx context.Context, info cameraservices.RequestInfo, event event) event {
	ctx, cancel := context.WithTimeout(ctx, 750*time.Second)
	defer cancel()

	// lookup hostname for source
	if info.SourceIP != nil {
		sources, _ := p.Resolver.LookupAddr(ctx, info.SourceIP.String())
		if len(sources) > 0 {
			event.User = sources[0]
		}
	}

	// lookup camera ip for building/room/device info
	if info.CameraIP != nil {
		event.TargetDevice.DeviceID = info.CameraIP.String()

		cameras, err := p.Resolver.LookupAddr(ctx, info.CameraIP.String())
		if err == nil && len(cameras) > 0 {
			trimmed := strings.TrimSuffix(cameras[0], ".byu.edu.")
			split := strings.SplitN(trimmed, "-", 3)
			if len(split) == 3 {
				event.TargetDevice.BuildingID = split[0]
				event.TargetDevice.RoomID = event.TargetDevice.BuildingID + "-" + split[1]
				event.TargetDevice.DeviceID = event.TargetDevice.RoomID + "-" + split[2]

				event.AffectedRoom = event.TargetDevice.roomInfo
			}
		}
	}

	return event
}

func (p *Publisher) publish(ctx context.Context, event event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("unable to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.URL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("unable to build request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("got a %v response from event url", resp.StatusCode)
	}

	return nil
}
