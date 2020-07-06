package cameraservices

import (
	"context"
)

type NewCameraFunc func(context.Context, string) (Camera, error)

type Camera interface {
	TiltUp(context.Context, byte) error
	TiltDown(context.Context, byte) error
	PanLeft(context.Context, byte) error
	PanRight(context.Context, byte) error
	PanTiltStop(context.Context) error
	ZoomTele(context.Context) error
	ZoomWide(context.Context) error
	ZoomStop(context.Context) error
	MemoryRecall(context.Context, byte) error
}
