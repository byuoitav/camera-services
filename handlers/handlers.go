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
		ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
		defer cancel()

		addrs, err := h.Resolver.LookupHost(ctx, addr)
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

func (h *Handlers) Lookup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	camIP, err := h.getCameraIP(ctx, c.Param("address"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(http.StatusOK, camIP.String())
}
