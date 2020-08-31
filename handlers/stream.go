package handlers

import (
	"bytes"
	"context"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

	m := multipart.NewWriter(c.Writer)

	// write the headers
	c.Writer.Header().Set(_hContentType, "multipart/x-mixed-replace; boundary="+m.Boundary())

	defer m.Close()

	frames := 0
	numErrs := 0
	avgFrameSize := 0
	start := time.Now()

	defer func() {
		avgFps := float64(frames) / time.Since(start).Seconds()
		log.Info("Done streaming", zap.Float64("avgFps", avgFps), zap.Int("avgFrameSize", avgFrameSize))
	}()

	buf := &bytes.Buffer{}
	header := textproto.MIMEHeader{}

	for {
		select {
		case image := <-images:
			buf.Reset()

			if err := jpeg.Encode(buf, image, nil); err != nil {
				log.Warn("unable to encode image", zap.Error(err))
				return
			}

			header.Set(_hContentType, "image/jpeg")
			header.Set(_hContentLength, strconv.Itoa(buf.Len()))

			part, err := m.CreatePart(header)
			if err != nil {
				log.Warn("unable to create part", zap.Error(err))
				return
			}

			if _, err := part.Write(buf.Bytes()); err != nil {
				log.Warn("unable to write part", zap.Error(err))
				return
			}

			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}

			if avgFrameSize == 0 {
				avgFrameSize = buf.Len()
			} else {
				avgFrameSize = (avgFrameSize + buf.Len()) / 2
			}

			numErrs = 0
			frames++
		case err := <-errs:
			numErrs++
			log.Warn("unable to get the next image", zap.Error(err))

			// we are averaging 8fps or so, so if we don't
			// get a single frame for 3 seconds, we'll end it
			if numErrs >= 24 {
				log.Warn("ending stream", zap.String("reason", "exceeded consecutive error count"))
				return
			}
		case <-ctx.Done():
			log.Info("ending stream", zap.String("reason", (ctx.Err().Error())))
			return
		}
	}
}
