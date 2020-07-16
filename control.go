package cameraservices

import "context"

type ControlKeyService interface {
	RoomAndControlGroup(ctx context.Context, key string) (string, string, error)
}

type ConfigService interface {
	Cameras(ctx context.Context, room string) ([]CameraConfig, error)
}

type CameraConfig struct {
	DisplayName string `json:"displayName"`

	TiltUp      string `json:"tiltUp"`
	TiltDown    string `json:"tiltDown"`
	PanLeft     string `json:"panLeft"`
	PanRight    string `json:"panRight"`
	PanTiltStop string `json:"panTiltStop"`

	ZoomIn   string `json:"zoomIn"`
	ZoomOut  string `json:"zoomOut"`
	ZoomStop string `json:"zoomStop"`

	Stream string `json:"stream"`

	Presets []CameraPreset `json:"presets"`
}

type CameraPreset struct {
	DisplayName string `json:"displayName"`
	SetPreset   string `json:"setPreset"`
}
