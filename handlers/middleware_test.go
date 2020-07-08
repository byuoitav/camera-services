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

func TestRequestIDWithID(t *testing.T) {
	id := "ID"

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	c.Request.Header.Set(_hRequestID, id)

	handlers := &Handlers{}
	handlers.RequestID(c)

	newID := c.GetString(_cRequestID)
	if newID != id {
		t.Fatalf("expected %q as request id, got %q", id, newID)
	}
}

func TestRequestIDWithoutID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handlers := &Handlers{}
	handlers.RequestID(c)

	id := c.GetString(_cRequestID)
	if id == "" {
		t.Fatalf("expected a random request id, but didn't get one")
	}
}

func TestCameraNoAddr(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	handler := &Handlers{}
	handler.Camera(c)

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

	handler := &Handlers{
		Logger:       log,
		CreateCamera: create,
	}
	handler.Camera(c)

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

	handler := &Handlers{
		Logger:       log,
		CreateCamera: create,
	}
	handler.Camera(c)

	if _, exists := c.Get(_cCamera); !exists {
		t.Fatalf("no camera found")
	}
}

func TestLog(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := &Handlers{
		Logger: log,
	}

	handler.Log(c)
}
