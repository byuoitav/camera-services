package handlers

import (
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
	_hRequestID     = "X-Request-ID"
	_hContentType   = "Content-Type"
	_hContentLength = "Content-Length"
)

type Middleware struct {
	Logger *zap.Logger
}

func (m *Middleware) RequestID(c *gin.Context) {
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

func (m *Middleware) Log(c *gin.Context) {
	id := c.GetString(_cRequestID)
	log := m.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	start := time.Now()
	log.Info("Starting request", zap.String("from", c.ClientIP()), zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path))
	c.Next()

	log.Info("Finished request", zap.Int("statusCode", c.Writer.Status()), zap.Duration("took", time.Since(start)))
}
