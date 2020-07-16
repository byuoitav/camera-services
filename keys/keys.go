package keys

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ControlKeyService struct {
	Address string
}

type roomControlGroupResponse struct {
	Room         string `json:"RoomID"`
	ControlGroup string `json:"PresetName"`
}

func (c *ControlKeyService) RoomAndControlGroup(ctx context.Context, key string) (string, string, error) {
	url := fmt.Sprintf("http://%s/%s/getPreset", c.Address, key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("unable to build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("unable to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return "", "", fmt.Errorf("invalid control key")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("unable to read response: %w", err)
	}

	var room roomControlGroupResponse
	if err := json.Unmarshal(body, &room); err != nil {
		return "", "", fmt.Errorf("unable to parse response: %w", err)
	}

	return room.Room, room.ControlGroup, nil
}
