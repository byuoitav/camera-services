package handlers

import (
	"bytes"
	"context"
	"fmt"
	"image/jpeg"
	"net/http"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	_mjpegBoundary     = "mjpeg_boundary"
	_mjpegFrameHeaderf = "\r\n--" + _mjpegBoundary + "\r\nContent-Type: image/jpeg\r\nContent-Length: %d\r\n\r\n"
)

func (h *Handlers) Stream(c *gin.Context) {
	cam := c.MustGet(_cCamera).(cameraservices.Camera)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Starting a stream")

	images, errs, err := cam.Stream(ctx)
	if err != nil {
		log.Warn("unable to start stream", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// write the normal http headers
	c.Writer.Header().Add(_hContentType, "multipart/x-mixed-replace;boundary="+_mjpegBoundary)

	defer c.Writer.WriteString("\r\n--" + _mjpegBoundary + "--")
	defer log.Info("Done streaming")

	buf := &bytes.Buffer{}
	for {
		select {
		case image := <-images:
			buf.Reset()

			if err := jpeg.Encode(buf, image, nil); err != nil {
				log.Warn("unable to encode image", zap.Error(err))
				continue
			}

			// write header for this frame
			c.Writer.WriteString(fmt.Sprintf(_mjpegFrameHeaderf, buf.Len()))
			c.Writer.Write(buf.Bytes())
		case err := <-errs:
			log.Warn("unable to get the next image", zap.Error(err))
			return
		case <-ctx.Done():
			log.Info("Ending stream", zap.String("reason", (ctx.Err().Error())))
			return
		}
	}
}
