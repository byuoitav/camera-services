package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGoToPresetFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	cam := badTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")
	c.Params = gin.Params{
		{
			Key:   "preset",
			Value: "1",
		},
	}

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.GoToPreset(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "go to preset error" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestMemoryRecallPass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	cam := goodTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")
	c.Params = gin.Params{
		{
			Key:   "channel",
			Value: "1",
		},
	}

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.GoToPreset(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}
