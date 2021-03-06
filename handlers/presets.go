package handlers

import (
	"context"
	"net/http"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *CameraController) GoToPreset(c *gin.Context) {
	cam := c.MustGet(_cCamera).(cameraservices.Camera)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	preset := c.Param("preset")
	log.Info("Going to preset", zap.String("preset", preset))

	if err := cam.GoToPreset(ctx, preset); err != nil {
		log.Warn("unable to go to preset", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Went to preset")
	c.Status(http.StatusOK)
}

func (h *CameraController) SavePreset(c *gin.Context) {
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

	preset := c.Param("preset")
	log.Info("Setting preset", zap.String("preset", preset))

	if err := cam.SetPreset(ctx, preset); err != nil {
		log.Warn("unable to set preset", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Preset set")
	c.Status(http.StatusOK)
}
