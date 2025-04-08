package couch

import (
	"context"
	"fmt"
	"strings"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/go-kivik/kivik/v3"
)

type configService struct {
	client     *kivik.Client
	uiConfigDB string
}

// New creates a new ConfigService, created a couchdb client pointed at url.
func New(ctx context.Context, url string, opts ...Option) (*configService, error) {
	client, err := kivik.New("couch", url)
	if err != nil {
		return nil, fmt.Errorf("unable to build client: %w", err)
	}

	return NewWithClient(ctx, client, opts...)
}

// NewWithClient creates a new ConfigService using the given client.
func NewWithClient(ctx context.Context, client *kivik.Client, opts ...Option) (*configService, error) {
	options := options{
		uiConfigDB: _defaultUIConfigDB,
	}

	for _, o := range opts {
		o.apply(&options)
	}

	if options.authFunc != nil {
		if err := client.Authenticate(ctx, options.authFunc); err != nil {
			return nil, fmt.Errorf("unable to authenticate: %w", err)
		}
	}

	return &configService{
		client:     client,
		uiConfigDB: options.uiConfigDB,
	}, nil
}

func (c *configService) Cameras(ctx context.Context, info cameraservices.ControlInfo) ([]cameraservices.CameraConfig, error) {
	var config uiConfig

	db := c.client.DB(ctx, c.uiConfigDB)
	if err := db.Get(ctx, info.Room).ScanDoc(&config); err != nil {
		return []cameraservices.CameraConfig{}, fmt.Errorf("unable to get/scan ui config: %w", err)
	}

	for _, cg := range config.ControlGroups {
		if cg.ID == info.ControlGroup && len(cg.Cameras) > 0 {
			return cg.Cameras, nil
		}
	}

	return []cameraservices.CameraConfig{}, fmt.Errorf("no cameras found in %s/%s", info.Room, info.ControlGroup)
}

func (c *configService) CameraPreset(ctx context.Context, camID, presetID string) (string, error) {
	db := c.client.DB(ctx, c.uiConfigDB)
	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"presets": map[string]interface{}{
				"$elemMatch": map[string]interface{}{
					"cameras": map[string]interface{}{
						"$elemMatch": map[string]interface{}{
							"presets": map[string]interface{}{
								"$elemMatch": map[string]interface{}{
									"setPreset": map[string]interface{}{
										"$regex": fmt.Sprintf(".*%s.*", camID),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	rows, err := db.Find(ctx, query)
	if err != nil {
		return "", fmt.Errorf("unable to find: %w", err)
	}

	if !rows.Next() {
		return "", fmt.Errorf("no matching documents found")
	}

	var config uiConfig
	if err := rows.ScanDoc(&config); err != nil {
		return "", fmt.Errorf("unable to scan doc: %w", err)
	}

	for _, cg := range config.ControlGroups {
		for _, cam := range cg.Cameras {
			for _, preset := range cam.Presets {
				if strings.Contains(preset.SetPreset, camID) && strings.HasSuffix(preset.SetPreset, presetID) {
					return preset.DisplayName, nil
				}
			}
		}
	}

	return "", fmt.Errorf("unable to find matching preset")
}

func (c *configService) Rooms(ctx context.Context) ([]string, error) {
	var rooms []string
	db := c.client.DB(ctx, c.uiConfigDB)

	query := map[string]interface{}{
		"fields": []string{"_id"},
		"limit":  2048,
		"selector": map[string]interface{}{
			"presets": map[string]interface{}{
				"$elemMatch": map[string]interface{}{
					"cameras": map[string]interface{}{
						"$elemMatch": map[string]interface{}{
							"displayName": map[string]interface{}{
								"$regex": "..*",
							},
						},
					},
				},
			},
		},
	}

	rows, err := db.Find(ctx, query)
	if err != nil {
		return rooms, fmt.Errorf("unable to find: %w", err)
	}

	for rows.Next() {
		var config uiConfig
		if err := rows.ScanDoc(&config); err != nil {
			continue
		}

		// rows.ID() should work, but i can't figure out why it doesn't...
		// rooms = append(rooms, rows.ID())
		rooms = append(rooms, config.ID)
	}

	return rooms, nil
}

func (c *configService) ControlGroups(ctx context.Context, room string) ([]string, error) {
	var config uiConfig
	var groups []string

	db := c.client.DB(ctx, c.uiConfigDB)
	if err := db.Get(ctx, room).ScanDoc(&config); err != nil {
		return groups, fmt.Errorf("unable to get/scan ui config: %w", err)
	}

	for _, cg := range config.ControlGroups {
		if len(cg.Cameras) > 0 {
			groups = append(groups, cg.ID)
		}
	}

	return groups, nil
}

// Returns a list of the urls for the commands that each contain the IP address or hostname
func (c *configService) ControlIP(ctx context.Context, room string) ([]string, error) {
	var config uiConfig
	var IP []string

	db := c.client.DB(ctx, c.uiConfigDB)
	if err := db.Get(ctx, room).ScanDoc(&config); err != nil {
		return IP, fmt.Errorf("unable to get/scan ui config: %w", err)
	}

	for _, cg := range config.ControlGroups {
		if len(cg.Cameras) > 0 {
			for _, cam := range cg.Cameras {
				IP = append(IP, cam.Stream)
			}
		}
	}

	return IP, nil
}
