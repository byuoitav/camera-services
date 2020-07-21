package handlers

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
)

func TestGetCameraIP(t *testing.T) {
	handler := CameraController{}
	ip, err := handler.getCameraIP(context.TODO(), "10:10")
	if err != nil {
		t.Fatalf("returned an error: %s", err)
	}

	if ip.String() != "0.0.0.10" {
		t.Fatalf("wrong ip returned")
	}
}

func TestCameraNoAddr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handler := &CameraController{}
	handler.CameraMiddleware(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "must include camera address" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestCameraFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Set(_cRequestID, "ID")
	c.Params = gin.Params{
		{
			Key:   "address",
			Value: "address",
		},
	}

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	create := func(ctx context.Context, addr string) (cameraservices.Camera, error) {
		return nil, errors.New("couldn't create camera")
	}

	handler := &CameraController{
		Logger:       log,
		CreateCamera: create,
	}
	handler.CameraMiddleware(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "unable to create camera couldn't create camera" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestCameraPass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		{
			Key:   "address",
			Value: "address",
		},
	}

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	create := func(ctx context.Context, addr string) (cameraservices.Camera, error) {
		return &goodTestCamera{}, nil
	}

	handler := &CameraController{
		Logger:       log,
		CreateCamera: create,
	}
	handler.CameraMiddleware(c)

	if _, exists := c.Get(_cCamera); !exists {
		t.Fatalf("no camera found")
	}
}
