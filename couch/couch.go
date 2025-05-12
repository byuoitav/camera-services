package couch

import (
	"context"
	"fmt"
	"net"
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

// FindCameraByAddress searches the config database for a camera with the given address
// and returns its username and password.
func (c *configService) FindCameraAuthByAddress(ctx context.Context, addr string) (string, string, error) {
	db := c.client.DB(ctx, c.uiConfigDB)

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"presets": map[string]interface{}{
				"$elemMatch": map[string]interface{}{
					"cameras": map[string]interface{}{
						"$elemMatch": map[string]interface{}{
							"$or": []map[string]interface{}{
								{"stream": map[string]interface{}{"$regex": addr}},
								{"panLeft": map[string]interface{}{"$regex": addr}},
								{"panRight": map[string]interface{}{"$regex": addr}},
								{"tiltUp": map[string]interface{}{"$regex": addr}},
								{"tiltDown": map[string]interface{}{"$regex": addr}},
							},
						},
					},
				},
			},
		},
	}

	rows, err := db.Find(ctx, query)
	if err != nil {
		return "", "", fmt.Errorf("unable to query config DB: %w", err)
	}
	for rows.Next() {
		var config uiConfig
		if err := rows.ScanDoc(&config); err != nil {
			continue
		}
		for _, group := range config.ControlGroups {
			for _, cam := range group.Cameras {
				if strings.Contains(cam.Stream, addr) || strings.Contains(cam.PanLeft, addr) ||
					strings.Contains(cam.PanRight, addr) || strings.Contains(cam.TiltUp, addr) ||
					strings.Contains(cam.TiltDown, addr) {
					return cam.UserName, cam.Password, nil
				}
			}
		}
	}
	return "", "", fmt.Errorf("camera with address %s not found", addr)
}

// GetCameraAuth retrieves the camera's username and password from the config database based on the provided address.
func (c *configService) GetCameraAuth(ctx context.Context, addr string) (string, string, error) {
	isIP := net.ParseIP(addr) != nil
	if !isIP {
		// Tries to get room from hostname in address
		parts := strings.SplitN(addr, "-", 3)
		if len(parts) >= 2 {
			roomID := parts[0] + "-" + parts[1]
			var config uiConfig

			db := c.client.DB(ctx, c.uiConfigDB)
			if err := db.Get(ctx, roomID).ScanDoc(&config); err == nil {
				for _, group := range config.ControlGroups {
					for _, cam := range group.Cameras {
						if cameraMatchesAddress(cam, addr) {
							return cam.UserName, cam.Password, nil
						}
					}
				}
			}
		}
	}

	// If it fails to get a room from the address, it will try to get the camera from the config database
	// by searching all the docs for the address
	return c.FindCameraAuthByAddress(ctx, addr)
}

func cameraMatchesAddress(cam cameraservices.CameraConfig, addr string) bool {
	return strings.Contains(cam.Stream, addr) ||
		strings.Contains(cam.PanLeft, addr) ||
		strings.Contains(cam.PanRight, addr) ||
		strings.Contains(cam.TiltUp, addr) ||
		strings.Contains(cam.TiltDown, addr)
}
