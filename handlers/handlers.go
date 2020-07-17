package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handlers struct {
	CreateCamera   cameraservices.NewCameraFunc
	EventPublisher cameraservices.EventPublisher
	Resolver       *net.Resolver
	Logger         *zap.Logger
}

func (h *Handlers) getCameraIP(ctx context.Context, addr string) (net.IP, error) {
	host := addr
	var err error

	if strings.Contains(host, ":") {
		host, _, err = net.SplitHostPort(host)
		if err != nil {
			return nil, fmt.Errorf("unable to split host/port: %w", err)
		}
	}

	// figure out if it's an ip or not
	ip := net.ParseIP(host)
	if ip == nil {
		addrs, err := h.Resolver.LookupHost(ctx, host)
		if err != nil {
			return nil, fmt.Errorf("unable to reverse lookup ip: %w", err)
		}

		if len(addrs) == 0 {
			return nil, errors.New("no camera IP addresses found")
		}

		for _, addr := range addrs {
			ip = net.ParseIP(addr)
			if ip != nil {
				break
			}
		}
	}

	return ip, nil
}

type ControlHandlers struct {
	ConfigService     cameraservices.ConfigService
	ControlKeyService cameraservices.ControlKeyService
}

func (h *ControlHandlers) GetCameras(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	room, _, err := h.ControlKeyService.RoomAndControlGroup(ctx, c.Param("key"))
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get room and control group: %s", err))
		return
	}

	cameras, err := h.ConfigService.Cameras(ctx, room)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get cameras: %s", err))
		return
	}

	c.JSON(http.StatusOK, cameras)
}
