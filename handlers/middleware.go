package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"go.uber.org/zap"
)

const (
	_cRequestID = "requestID"
	_cCamera    = "camera"
)

const (
	_hRequestID   = "X-Request-ID"
	_hContentType = "Content-Type"
)

func (h *Handlers) RequestID(c *gin.Context) {
	var id string
	if c.GetHeader(_hRequestID) != "" {
		id = c.GetHeader(_hRequestID)
	} else {
		uid, err := ksuid.NewRandom()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		id = uid.String()
	}

	c.Set(_cRequestID, id)
	c.Next()
}

func (h *Handlers) Camera(c *gin.Context) {
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

func (h *Handlers) Log(c *gin.Context) {
	id := c.GetString(_cRequestID)
	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	start := time.Now()
	log.Info("Starting request", zap.String("from", c.ClientIP()), zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path))
	c.Next()

	log.Info("Finished request", zap.Int("statusCode", c.Writer.Status()), zap.Duration("took", time.Since(start)))
}
