package handlers

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"sync"
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

type stream struct {
	subs  int
	frame []byte
	sync.RWMutex
}

func (h *CameraController) stream(c *gin.Context) {
	cam := c.MustGet(_cCamera).(cameraservices.Camera)
	id := c.GetString(_cRequestID)

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	// get the stream, or start it
	v, err, _ := h.single.Do("Stream+address", func() (interface{}, error) {
		s, ok := h.streams.Load(cam)
		if ok {
			return s, nil
		}

		s, err := h.startStream(cam)
		if err != nil {
			return nil, err
		}

		h.streams.Store(cam, s)
		return s, nil
	})
	if err != nil {
		// unable to start stream
		return
	}

	s := v.(*stream)
	s.Lock()
	s.subs++
	s.Unlock()

	// start a multipart writer
	m := multipart.NewWriter(c.Writer)
	c.Writer.Header().Set(_hContentType, "multipart/x-mixed-replace; boundary="+m.Boundary())
	defer m.Close()

	header := textproto.MIMEHeader{}
	header.Set(_hContentType, "image/jpeg")

	ticker := time.NewTicker(125 * time.Millisecond)
	defer ticker.Stop()

	// prevFrame to not rewrite a duplicate??

	for {
		select {
		case <-ticker.C:
			s.RLock()
			frame := make([]byte, len(s.frame))
			copy(frame, s.frame)
			s.RUnlock()

			header.Set(_hContentLength, strconv.Itoa(len(frame)))

			part, err := m.CreatePart(header)
			if err != nil {
				log.Warn("unable to create part", zap.Error(err))
				return
			}

			if _, err := part.Write(frame); err != nil {
				log.Warn("unable to write part", zap.Error(err))
				return
			}

			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-(make(chan string)): // just for lint rn, should be ctx probably
		}
	}
}

func (h *CameraController) startStream(cam cameraservices.Camera, s *stream) error {
	var jpegs chan []byte
	var errors chan error
	var err error

	if jCam, ok := cam.(cameraservices.JPEGCamera); ok {
		jpegs, errors, err = jCam.StreamJPEG(context.Background())
	} else {
		var imgs chan image.Image
		var streamErrs chan error
		convertErrs := make(chan error)
		jpegs = make(chan []byte)

		imgs, streamErrs, err = cam.Stream(context.Background())
		if err == nil {
			// convert imgs to jpegs
			go func() {
				defer close(jpegs)
				defer close(convertErrs)

				buf := &bytes.Buffer{}
				for img := range imgs {
					buf.Reset()

					if err := jpeg.Encode(buf, img, nil); err != nil {
						// log.Warn("unable to encode image", zap.Error(err))
						convertErrs <- err
						continue
					}

					jpegs <- buf.Bytes()
				}
			}()

			// merge streamErrs and convertErrs channels to errors
			go func() {
				defer close(errors)

				for {
					select {
					case err, ok := <-streamErrs:
						if !ok {
							return
						}

						errors <- err
					case err := <-convertErrs:
						if !ok {
							return
						}

						errors <- err
					}
				}
			}()
		}
	}

	if err != nil {
		// log.Warn("unable to start stream", zap.Error(err))
		// c.String(http.StatusInternalServerError, err.Error())
		return err
	}

	// check this often every if nobody is using our stream anymore
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// defer removing this stream from the map
	defer h.streams.Delete(cam)

	errCount := 0
	for {
		select {
		case jpeg, ok := <-jpegs:
			if !ok {
				return nil
			}

			s.Lock()
			s.frame = jpeg
			s.Unlock()
		case err, ok := <-errors:
			if !ok {
				return nil
			}

			errCount++
			if errCount >= 24 {
				return err
			}
		case <-ticker.C:
			s.RLock()
			if s.subs == 0 {
				s.RUnlock()
				return nil
			}

			s.RUnlock()
		}
	}
}
