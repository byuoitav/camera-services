package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDWithID(t *testing.T) {
	id := "ID"

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	c.Request.Header.Set(_hRequestID, id)

	m := &Middleware{}
	m.RequestID(c)

	newID := c.GetString(_cRequestID)
	if newID != id {
		t.Fatalf("expected %q as request id, got %q", id, newID)
	}
}

func TestRequestIDWithoutID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)

	m := &Middleware{}
	m.RequestID(c)

	id := c.GetString(_cRequestID)
	if id == "" {
		t.Fatalf("expected a random request id, but didn't get one")
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

	m := &Middleware{
		Logger: log,
	}

	m.Log(c)
}
