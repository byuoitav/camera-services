package handlers

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type CameraController struct {
	CreateCamera   cameraservices.NewCameraFunc
	EventPublisher cameraservices.EventPublisher
	Logger         *zap.Logger

	streams *sync.Map
	single  *singleflight.Group
}

func NewCameraController() *CameraController {
	return &CameraController{
		streams: &sync.Map{},
		single:  &singleflight.Group{},
	}
}

func (h *CameraController) getCameraIP(ctx context.Context, addr string) (net.IP, error) {
	var err error

	if strings.Contains(addr, ":") {
		addr, _, err = net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("unable to split host/port: %w", err)
		}
	}

	// figure out if it's an ip or not
	ip := net.ParseIP(addr)
	if ip == nil {
		ctx, cancel := context.WithTimeout(ctx, 750*time.Millisecond)
		defer cancel()

		addrs, err := net.DefaultResolver.LookupHost(ctx, addr)
		if err != nil {
			return nil, fmt.Errorf("unable to reverse lookup ip: %w", err)
		}

		if len(addrs) == 0 {
			return nil, errors.New("no camera IP addresses found")
		}

		for i := range addrs {
			ip = net.ParseIP(addrs[i])
			if ip != nil {
				break
			}
		}
	}

	return ip, nil
}

func (h *CameraController) CameraMiddleware(c *gin.Context) {
	addr := c.Param("address")
	if addr == "" {
		c.String(http.StatusBadRequest, "must include camera address")
		c.Abort()
		return
	}

	id := c.GetString(_cRequestID)
	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Debug("Getting camera", zap.String("address", addr))

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	cam, err := h.CreateCamera(ctx, addr)
	if err != nil {
		c.String(http.StatusInternalServerError, "unable to create camera %s", err)
		c.Abort()
		return
	}

	log.Debug("Got camera")

	c.Set(_cCamera, cam)
	c.Next()
}

func (h *CameraController) Reboot(c *gin.Context) {
	id := c.GetString(_cRequestID)
	cam, ok := c.MustGet(_cCamera).(cameraservices.CameraAdmin)
	if !ok || cam == nil {
		c.String(http.StatusBadRequest, "not supported")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Rebooting...")

	if err := cam.Reboot(ctx); err != nil {
		log.Warn("unable to reboot", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Rebooted")
	c.Status(http.StatusOK)
}
