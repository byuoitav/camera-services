package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMemoryRecallNoChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	cam := goodTestCamera{}
	c.Set(_cRequestID, "ID")
	c.Set(_cCamera, &cam)
	c.Params = gin.Params{
		{
			Key:   "channel",
			Value: "channel",
		},
	}

	log, err := SetLogger()
	if err != nil {
		t.Fatalf("unable to build logger: %s", err)
	}

	handler := Handlers{
		Logger: log,
	}

	handler.MemoryRecall(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "invalid channel: strconv.Atoi: parsing \"channel\": invalid syntax" {
		t.Fatalf("incorrect error generated: %s", string(body))
	}

}

func TestMemoryRecallFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	cam := badTestCamera{}
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

	handler := Handlers{
		Logger: log,
	}

	handler.MemoryRecall(c)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body: %s", err)
	}

	if string(body) != "memory recall error" {
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

	handler := Handlers{
		Logger: log,
	}

	handler.MemoryRecall(c)
	if resp.Result().StatusCode/100 != 2 {
		t.Fatalf("wrong response status code received: %d", resp.Result().StatusCode)
	}
}
