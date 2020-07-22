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

func (h *CameraController) Stream(c *gin.Context) {
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

	defer func() {
		_, _ = c.Writer.WriteString("\r\n--" + _mjpegBoundary + "--")
	}()

	frames := 0
	numErrs := 0
	start := time.Now()

	defer log.Info("Done streaming", zap.Float64("avgFps", float64(frames)/time.Since(start).Seconds()))

	buf := &bytes.Buffer{}
	for {
		select {
		case image := <-images:
			buf.Reset()
			numErrs = 0

			if err := jpeg.Encode(buf, image, nil); err != nil {
				log.Warn("unable to encode image", zap.Error(err))
				continue
			}

			// write header for this frame
			if _, err := c.Writer.WriteString(fmt.Sprintf(_mjpegFrameHeaderf, buf.Len())); err != nil {
				log.Warn("unable to write frame header", zap.Error(err))
				return
			}

			if _, err := c.Writer.Write(buf.Bytes()); err != nil {
				log.Warn("unable to write frame", zap.Error(err))
				return
			}

			frames++
		case err := <-errs:
			numErrs++
			log.Warn("unable to get the next image", zap.Error(err))

			if numErrs >= 3 {
				log.Warn("Ending stream", zap.String("reason", "exceeded consecutive error count"))
				return
			}
		case <-ctx.Done():
			log.Info("Ending stream", zap.String("reason", (ctx.Err().Error())))
			return
		}
	}
}
