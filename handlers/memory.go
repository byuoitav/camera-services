package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handlers) MemoryRecall(c *gin.Context) {
	cam := c.MustGet(_cCamera).(cameraservices.Camera)
	id := c.GetString(_cRequestID)

	channel, err := strconv.Atoi(c.Param("channel"))
	if err != nil {
		c.String(http.StatusBadRequest, "invalid channel: %s", err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Recalling memory", zap.Int("channel", channel))

	if err := cam.MemoryRecall(ctx, byte(channel)); err != nil {
		log.Warn("unable to recall memory", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Info("Recalled memory")
	c.Status(http.StatusOK)
}
