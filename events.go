package cameraservices

import (
	"context"
	"net"
	"time"
)

type EventPublisher interface {
	Publish(context.Context, RequestInfo) error
	Error(context.Context, RequestError) error
}

type RequestInfo struct {
	Action    string
	Timestamp time.Time
	SourceIP  net.IP
	CameraIP  net.IP
	Duration  time.Duration
	Data      map[string]interface{}
}

type RequestError struct {
	RequestInfo
	Error string
}
