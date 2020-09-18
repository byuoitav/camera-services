package cameraservices

import (
	"context"
	"image"
)

type NewCameraFunc func(context.Context, string) (Camera, error)

type Camera interface {
	TiltUp(context.Context) error
	TiltDown(context.Context) error
	PanLeft(context.Context) error
	PanRight(context.Context) error
	PanTiltStop(context.Context) error
	ZoomIn(context.Context) error
	ZoomOut(context.Context) error
	ZoomStop(context.Context) error
	GoToPreset(context.Context, string) error
	Stream(context.Context) (chan image.Image, chan error, error)
}

type Rebootable interface {
	Reboot(context.Context) error
}
