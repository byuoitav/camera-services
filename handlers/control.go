package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ControlHandlers struct {
	ConfigService     cameraservices.ConfigService
	ControlKeyService cameraservices.ControlKeyService
	Me                *url.URL
	Logger            *zap.Logger
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

	// change urls to go through proxy (me)
	rewrite := func(u string) string {
		url, err := url.Parse(u)
		if err != nil {
			return ""
		}

		url.Scheme = h.Me.Scheme
		url.Host = h.Me.Host

		switch {
		case strings.Contains(u, "aver"):
			url.Path = "/proxy/aver" + url.Path
		case strings.Contains(u, "axis"):
			url.Path = "/proxy/axis" + url.Path
		}

		return url.String()
	}

	for i := range cameras {
		cameras[i].PanLeft = rewrite(cameras[i].PanLeft)
		cameras[i].PanRight = rewrite(cameras[i].PanRight)
		cameras[i].TiltUp = rewrite(cameras[i].TiltUp)
		cameras[i].TiltDown = rewrite(cameras[i].TiltDown)
		cameras[i].PanTiltStop = rewrite(cameras[i].PanTiltStop)
		cameras[i].ZoomIn = rewrite(cameras[i].ZoomIn)
		cameras[i].ZoomOut = rewrite(cameras[i].ZoomOut)
		cameras[i].ZoomStop = rewrite(cameras[i].ZoomStop)
		cameras[i].Stream = rewrite(cameras[i].Stream)

		for j := range cameras[i].Presets {
			cameras[i].Presets[j].SetPreset = rewrite(cameras[i].Presets[j].SetPreset)
		}
	}

	c.JSON(http.StatusOK, cameras)
}

func (h *ControlHandlers) Proxy(to *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString(_cRequestID)

		log := h.Logger
		if len(id) > 0 {
			log = log.With(zap.String("requestID", id))
		}

		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					return
				}

				panic(err)
			}
		}()

		proxy := httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = to.Scheme
				req.URL.Host = to.Host

				switch {
				case strings.HasPrefix(req.URL.Path, "/proxy/aver"):
					req.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/proxy/aver")
				case strings.HasPrefix(req.URL.Path, "/proxy/axis"):
					req.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/proxy/axis")
				}

				req.Header.Set(_hRequestID, id)

				log.Debug("Forwarding request to", zap.String("url", req.URL.String()))
			},
			ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
				log.Warn("error proxying request", zap.Error(err))
				rw.WriteHeader(http.StatusBadGateway)
				_, _ = rw.Write([]byte(fmt.Sprintf("unable to proxy request: %s", err)))
			},
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
