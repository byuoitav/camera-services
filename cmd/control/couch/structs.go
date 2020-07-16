package couch

import control "github.com/byuoitav/camera-services/cmd/control/data"

type uiConfig struct {
	ControlGroups []struct {
		ID      string           `json:"name"`
		Cameras []control.Camera `json:"cameras"`
	} `json:"presets"`
}
