package couch

import cameraservices "github.com/byuoitav/camera-services"

type uiConfig struct {
	ID            string `json:"_id"`
	ControlGroups []struct {
		ID      string                        `json:"name"`
		Cameras []cameraservices.CameraConfig `json:"cameras"`
	} `json:"presets"`
}
