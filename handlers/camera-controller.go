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
	CreateCamera      cameraservices.NewCameraFunc
	EventPublisher    cameraservices.EventPublisher
	ControlKeyService cameraservices.ControlKeyService
	DatabaseService   cameraservices.ConfigService
	Logger            *zap.Logger

	streams *sync.Map
	single  *singleflight.Group
}

func NewCameraController(cs cameraservices.ConfigService) *CameraController {
	return &CameraController{
		streams:         &sync.Map{},
		single:          &singleflight.Group{},
		DatabaseService: cs,
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

func (h *CameraController) checkControlKey(c *gin.Context, key, address string) bool {
	if key == "" {
		return false
	}

	room, _, err := h.ControlKeyService.RoomAndControlGroup(c, key)
	if err != nil {
		return false
	}

	IPAddresses, err := h.DatabaseService.ControlIP(c, room)
	if err != nil {
		return false
	}

	// Check if any IPAddresses contain address string
	fmt.Println("Checking address: " + address)
	for _, IP := range IPAddresses {
		ipStr := fmt.Sprintf("%v", IP) // Convert IP to string
		fmt.Println(ipStr)
		if strings.Contains(ipStr, address) {
			return true
		}
	}
	return true
}

func (h *CameraController) CameraMiddleware(c *gin.Context) {
	ck, err := c.Cookie("control-key")
	if err != nil {
		c.String(http.StatusUnauthorized, "no control key")
		c.Abort()
		return
	}

	address := c.Param("address")
	if address == "" {
		c.String(http.StatusBadRequest, "must include camera address")
		c.Abort()
		return
	}

	authorized := h.checkControlKey(c, ck, address)
	if !authorized {
		c.String(http.StatusForbidden, "Unauthorized Key")
		c.Abort()
		return
	}

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
