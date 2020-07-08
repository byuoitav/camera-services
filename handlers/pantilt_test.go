package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestTiltUpPass(t *testing.T) {
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

	handler := Handlers{
		Logger: log,
	}

	handler.TiltUp(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestTiltUpFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := Handlers{
		Logger: log,
	}

	handler.TiltUp(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "tilt up error" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestTiltDownPass(t *testing.T) {
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

	handler := Handlers{
		Logger: log,
	}

	handler.TiltDown(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestTiltDownFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := Handlers{
		Logger: log,
	}

	handler.TiltDown(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "tilt down error" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestPanLeftPass(t *testing.T) {
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

	handler := Handlers{
		Logger: log,
	}

	handler.PanLeft(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestPanLeftFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := Handlers{
		Logger: log,
	}

	handler.PanLeft(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "pan left error" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestPanRightPass(t *testing.T) {
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

	handler := Handlers{
		Logger: log,
	}

	handler.PanRight(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestPanRightFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := Handlers{
		Logger: log,
	}

	handler.PanRight(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "pan right error" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}

func TestPanTiltStopPass(t *testing.T) {
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

	handler := Handlers{
		Logger: log,
	}

	handler.PanTiltStop(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}

func TestPanTiltStopFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	cam := badTestCamera{}
	c.Set(_cCamera, &cam)
	c.Set(_cRequestID, "ID")

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := Handlers{
		Logger: log,
	}

	handler.PanTiltStop(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "pan tilt stop error" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}
}
