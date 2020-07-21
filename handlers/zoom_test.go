package handlers

import (
	"context"
	"errors"
	"image"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
)

type goodTestCamera struct{}
type badTestCamera struct{}

func (t *goodTestCamera) TiltUp(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) TiltDown(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) PanLeft(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) PanRight(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) PanTiltStop(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) ZoomIn(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) ZoomOut(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) ZoomStop(ctx context.Context) error {
	return nil
}

func (t *goodTestCamera) GoToPreset(ctx context.Context, preset string) error {
	return nil
}

// TODO test with this. just added so that the tests will build
func (t *goodTestCamera) Stream(ctx context.Context) (chan image.Image, chan error, error) {
	return nil, nil, nil
}

func (t *badTestCamera) TiltUp(ctx context.Context) error {
	return errors.New("tilt up error")
}

func (t *badTestCamera) TiltDown(ctx context.Context) error {
	return errors.New("tilt down error")
}

func (t *badTestCamera) PanLeft(ctx context.Context) error {
	return errors.New("pan left error")
}

func (t *badTestCamera) PanRight(ctx context.Context) error {
	return errors.New("pan right error")
}

func (t *badTestCamera) PanTiltStop(ctx context.Context) error {
	return errors.New("pan tilt stop error")
}

func (t *badTestCamera) ZoomIn(ctx context.Context) error {
	return errors.New("no zoom in")
}

func (t *badTestCamera) ZoomOut(ctx context.Context) error {
	return errors.New("no zoom out")
}

func (t *badTestCamera) ZoomStop(ctx context.Context) error {
	return errors.New("no zoom stop")
}

func (t *badTestCamera) GoToPreset(ctx context.Context, preset string) error {
	return errors.New("go to preset error")
}

// TODO test with this. just added so that the tests will build
func (t *badTestCamera) Stream(ctx context.Context) (chan image.Image, chan error, error) {
	return nil, nil, nil
}

func SetLogger() (*zap.Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(1),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json", EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}

func TestZoomInPass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := goodTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.ZoomIn(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestZoomInFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.ZoomIn(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "no zoom in" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestZoomOutPass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := goodTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.ZoomOut(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestZoomOutFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.ZoomOut(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "no zoom out" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestZoomStopPass(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := goodTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.ZoomStop(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestZoomStopFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := CameraController{
		Logger: log,
	}

	handler.ZoomStop(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "no zoom stop" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}
