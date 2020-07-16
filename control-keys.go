package cameraservices

import "context"

type ControlKeyService interface {
	RoomAndControlGroup(ctx context.Context, key string) (string, string, error)
}
