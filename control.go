package cameraservices

import "context"

type ControlKeyService interface {
	RoomAndControlGroup(ctx context.Context, key string) (string, string, error)
}

type ConfigService interface {
	Cameras(context.Context, ControlInfo) ([]CameraConfig, error)
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

	// admin items
	Reboot string `json:"reboot"`
}

type CameraPreset struct {
	DisplayName string `json:"displayName"`
	SavePreset  string `json:"savePreset"`
	SetPreset   string `json:"setPreset"`
}

type ControlInfo struct {
	Room         string `json:"room" form:"room"`
	ControlGroup string `json:"controlGroup" form:"controlGroup"`
	ControlKey   string `json:"controlKey" form:"controlKey"`
}
