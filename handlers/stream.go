package handlers

import (
	"context"
	"net/http"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	_mjpegBoundary = "MJPEG_BOUNDARY"
)

func (h *Handlers) Stream(c *gin.Context) {
	cam := c.MustGet(_cCamera).(cameraservices.Camera)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Minute)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Starting a stream")

	_, err := cam.Stream(ctx)
	if err != nil {
		log.Warn("unable to stream", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// add the stream headers
	c.Writer.Header().Add(_hContentType, "multipart/x-mixed-replace;boundary="+_mjpegBoundary)

	log.Info("Done streaming")
	c.Status(http.StatusOK)
}
