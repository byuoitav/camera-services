package couch

import (
	"context"
	"fmt"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/go-kivik/kivik"
)

type configService struct {
	client     *kivik.Client
	uiConfigDB string
}

// New creates a new ConfigService, created a couchdb client pointed at url.
func New(ctx context.Context, url string, opts ...Option) (cameraservices.ConfigService, error) {
	client, err := kivik.New("couch", url)
	if err != nil {
		return nil, fmt.Errorf("unable to build client: %w", err)
	}

	return NewWithClient(ctx, client, opts...)
}

// NewWithClient creates a new ConfigService using the given client.
func NewWithClient(ctx context.Context, client *kivik.Client, opts ...Option) (cameraservices.ConfigService, error) {
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

func (c *configService) Cameras(ctx context.Context, room string) ([]cameraservices.CameraConfig, error) {
	var config uiConfig

	db := c.client.DB(ctx, c.uiConfigDB)
	if err := db.Get(ctx, room).ScanDoc(&config); err != nil {
		return []cameraservices.CameraConfig{}, fmt.Errorf("unable to get/scan ui config: %w", err)
	}

	for _, cg := range config.ControlGroups {
		if cg.Cameras != nil && len(cg.Cameras) != 0 {
			return cg.Cameras, nil
		}
	}

	return []cameraservices.CameraConfig{}, fmt.Errorf("no cameras found in %s", room)
}
