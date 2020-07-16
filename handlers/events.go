package handlers

import (
	"bytes"
	"context"
	"net"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type responseSaver struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseSaver) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseSaver) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func (h *Handlers) Publish(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		info := cameraservices.RequestInfo{
			Action:    action + c.Param("channel"),
			Timestamp: time.Now(),
		}

		id := c.GetString(_cRequestID)
		log := h.Logger
		if len(id) > 0 {
			log = log.With(zap.String("requestID", id))
		}

		info.SourceIP = net.ParseIP(c.ClientIP())

		w := &responseSaver{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// run the rest of the handlers
		c.Next()

		go func(status int) {
			info.Duration = time.Since(info.Timestamp)

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			var err error
			info.CameraIP, err = h.getCameraIP(ctx, c.Param("address"))
			if err != nil {
				log.Warn("unable to get camera ip", zap.Error(err))
			}

			if status/100 != 2 {
				err = h.EventPublisher.Error(ctx, cameraservices.RequestError{
					RequestInfo: info,
					Error:       w.body.String(),
				})
			} else {
				err = h.EventPublisher.Publish(ctx, info)
			}

			if err != nil {
				log.Warn("unable to publish event", zap.Error(err))
			}
		}(c.Writer.Status())
	}
}
