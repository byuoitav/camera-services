package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-kivik/couchdb/v3"
	"github.com/go-kivik/kivik/v3"
)

type uiConfig struct {
	data map[string]interface{}

	Presets []preset `json:"presets"`
}

type preset struct {
	data map[string]interface{}

	Cameras []struct {
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

		Reboot string `json:"reboot"`

		Presets []struct {
			DisplayName string `json:"displayName"`
			SetPreset   string `json:"setPreset"`
			SavePreset  string `json:"savePreset"`
		} `json:"presets"`
	} `json:"cameras"`
}

type info struct {
	RoomID     string
	PresetName string
	Address    string
	Type       string
}

func (u uiConfig) MarshalJSON() ([]byte, error) {
	u.data["presets"] = u.Presets
	return json.Marshal(u.data)
}

func (u *uiConfig) UnmarshalJSON(data []byte) error {
	type Alias uiConfig
	cfg := &struct {
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}

	u.data = make(map[string]interface{})
	if err := json.Unmarshal(data, &u.data); err != nil {
		return err
	}

	return nil
}

func (p preset) MarshalJSON() ([]byte, error) {
	p.data["cameras"] = p.Cameras
	return json.Marshal(p.data)
}

func (p *preset) UnmarshalJSON(data []byte) error {
	type Alias preset
	cfg := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return err
	}

	p.data = make(map[string]interface{})

	if err := json.Unmarshal(data, &p.data); err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := kivik.New("couch", os.Getenv("DB_ADDRESS"))
	if err != nil {
		log.Fatalf("unable to build couch client: %s", err)
	}

	client.Authenticate(ctx, couchdb.BasicAuth(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD")))

	db := client.DB(ctx, "ui-configuration")
	rows, err := db.AllDocs(ctx, kivik.Options{"include_docs": true})
	if err != nil {
		log.Fatalf("unable to get all docs: %s", err)
	}

	// use csv
	fmt.Printf("room,preset,type,address,error\n")

	for rows.Next() {
		var cfg uiConfig
		if err := rows.ScanDoc(&cfg); err != nil {
			log.Fatalf("unable to scan %q: %s", rows.ID(), err)
		}

		if cfg.data["_id"] != "MARB-108" {
			continue
		}

		var changed bool
		var hasCamera bool
		var info info
		for i := range cfg.Presets {
			if len(cfg.Presets[i].Cameras) > 0 {
				hasCamera = true
				log.Printf("Found camera on %s/%s", rows.ID(), cfg.Presets[i].data["name"])
				info.RoomID = rows.ID()
				info.PresetName = cfg.Presets[i].data["name"].(string)

				// change all of the cameras
				for j := range cfg.Presets[i].Cameras {
					switch {
					case strings.Contains(cfg.Presets[i].Cameras[j].TiltUp, "aver"):
						changed = true
						addr := strings.TrimPrefix(cfg.Presets[i].Cameras[j].TiltUp, "https://aver.av.byu.edu/v1/Pro520/")
						addr = strings.TrimPrefix(addr, "https://aver-dev.av.byu.edu/v1/Pro520/")
						addr = strings.TrimSuffix(addr, "/pantilt/up")

						log.Printf("Changing %s to aver (addr=%q)", rows.ID(), addr)
						info.Type = "aver"
						info.Address = addr

						cfg.Presets[i].Cameras[j].Reboot = fmt.Sprintf("https://aver.av.byu.edu/v1/Pro520/%s/reboot", addr)
						for k := range cfg.Presets[i].Cameras[j].Presets {
							cfg.Presets[i].Cameras[j].Presets[k].SavePreset = fmt.Sprintf("https://aver.av.byu.edu/v1/Pro520/%s/savePreset/%d", addr, k)
						}
					default:
						info.Type = "unknown"
						log.Printf("Nothing to do for %s/%s", rows.ID(), cfg.Presets[i].Cameras[j].DisplayName)
					}
				}
			}
		}

		if !hasCamera {
			continue
		}

		if !changed {
			log.Printf("Not posting a new document for %s", rows.ID())
			fmt.Printf("%s,%s,%s,%s,%s\n", info.RoomID, info.PresetName, info.Type, info.Address, "")
			continue
		}

		if _, err := db.Put(ctx, rows.ID(), cfg); err != nil {
			fmt.Printf("%s,%s,%s,%s,%s\n", info.RoomID, info.PresetName, info.Type, info.Address, err.Error())
			log.Printf("failed to update %q: %s", rows.ID(), err)
			continue
		}

		fmt.Printf("%s,%s,%s,%s,%s\n", info.RoomID, info.PresetName, info.Type, info.Address, "success")
	}
}
