package opa

import (
	"context"
	"testing"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/stretchr/testify/require"
)

func TestAuthMap(t *testing.T) {
	ctx := context.Background()
	ctx = cameraservices.WithAuth(ctx, map[string]bool{
		"reboot": true,
	})

	client := &Client{}
	require.True(t, client.IsAuthorizedFor(ctx, "reboot"))
	require.False(t, client.IsAuthorizedFor(ctx, "setPreset"))
}
