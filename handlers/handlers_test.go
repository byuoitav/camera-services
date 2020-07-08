package handlers

import (
	"context"
	"testing"
)

func TestGetCameraIP(t *testing.T) {
	handler := Handlers{}
	ip, err := handler.getCameraIP(context.TODO(), "10:10")
	if err != nil {
		t.Fatalf("returned an error: %s", err)
	}

	if ip.String() != "0.0.0.10" {
		t.Fatalf("wrong ip returned")
	}
}
