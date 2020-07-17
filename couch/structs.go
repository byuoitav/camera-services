package couch

import cameraservices "github.com/byuoitav/camera-services"

type uiConfig struct {
	ControlGroups []struct {
		ID      string                        `json:"name"`
		Cameras []cameraservices.CameraConfig `json:"cameras"`
	} `json:"presets"`
}
