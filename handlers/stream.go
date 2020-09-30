package handlers

import (
	"bytes"
	"context"
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

type stream struct {
	sync.Mutex
	subs map[chan []byte]struct{}
	done chan struct{}
}

func (h *CameraController) Stream(c *gin.Context) {
	cam := c.MustGet(_cCamera).(cameraservices.Camera)
	id := c.GetString(_cRequestID)

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Minute)
	defer cancel()

	log := h.Logger
	if len(id) > 0 {
		log = log.With(zap.String("requestID", id))
	}

	log.Info("Subscribing to stream")

	// get the stream, or start it
	v, err, _ := h.single.Do("stream"+cam.RemoteAddr(), func() (interface{}, error) {
		s, ok := h.streams.Load(cam)
		if ok {
			return s, nil
		}

		s, err := h.startStream(cam, h.Logger.With(zap.String("addr", cam.RemoteAddr())))
		if err != nil {
			return nil, err
		}

		h.streams.Store(cam, s)
		return s, nil
	})
	if err != nil {
		log.Warn("unable to start stream", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// add our frames channel to the stream so they get sent to us
	frames := make(chan []byte)

	s := v.(*stream)
	s.Lock()
	s.subs[frames] = struct{}{}
	log.Info("Subscribed to stream", zap.Int("numSubs", len(s.subs)))
	s.Unlock()

	// start a multipart writer
	m := multipart.NewWriter(c.Writer)
	defer m.Close()

	// write the headers
	c.Writer.Header().Set(_hContentType, "multipart/x-mixed-replace; boundary="+m.Boundary())

	// headers for each frame
	header := textproto.MIMEHeader{}
	header.Set(_hContentType, "image/jpeg")

	defer func() {
		s.Lock()
		delete(s.subs, frames)
		log.Info("Unsubscribing from stream", zap.Int("numSubs", len(s.subs)))
		s.Unlock()
	}()

	for {
		select {
		case frame, ok := <-frames:
			if !ok {
				log.Info("Finished streaming", zap.String("reason", "frame chan closed"))
				return
			}

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

		case <-ctx.Done():
			log.Info("Finished streaming", zap.String("reason", ctx.Err().Error()))
			return
		case <-s.done:
			log.Info("Finished streaming", zap.String("reason", "done chan closed"))
			return
		}
	}
}

func (h *CameraController) startStream(cam cameraservices.Camera, log *zap.Logger) (*stream, error) {
	var jpegs chan []byte
	var errors chan error

	ctx, cancel := context.WithCancel(context.Background())
	log.Info("Starting stream")

	if jCam, ok := cam.(cameraservices.JPEGCamera); ok {
		var err error

		jpegs, errors, err = jCam.StreamJPEG(ctx)
		if err != nil {
			cancel()
			return nil, err
		}
	} else {
		imgs, streamErrs, err := cam.Stream(ctx)
		if err != nil {
			cancel()
			return nil, err
		}

		jpegs = make(chan []byte)
		errors = make(chan error)
		convertErrs := make(chan error)

		// convert imgs to jpegs
		go func() {
			defer close(jpegs)
			defer close(convertErrs)

			buf := &bytes.Buffer{}
			for img := range imgs {
				buf.Reset()

				if err := jpeg.Encode(buf, img, nil); err != nil {
					convertErrs <- err
					continue
				}

				select {
				case jpegs <- buf.Bytes():
				default:
				}
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

					select {
					case errors <- err:
					default:
					}
				case err, ok := <-convertErrs:
					if !ok {
						return
					}

					select {
					case errors <- err:
					default:
					}
				}
			}
		}()
	}

	log.Info("Started stream")

	s := &stream{
		subs: make(map[chan []byte]struct{}),
		done: make(chan struct{}),
	}

	go func() {
		// metrics info
		frameCount := 0
		avgFrameSize := 0
		start := time.Now()

		defer func() {
			h.streams.Delete(cam)
			close(s.done)
			cancel()

			avgFps := float64(frameCount) / time.Since(start).Seconds()
			log.Info("Stopped stream", zap.Float64("avgFps", avgFps), zap.Int("avgFrameSize", avgFrameSize), zap.Duration("duration", time.Since(start)))
		}()

		errCount := 0
		for {
			select {
			case jpeg, ok := <-jpegs:
				if !ok {
					return
				}

				s.Lock()
				if len(s.subs) == 0 {
					log.Info("No more subs on stream, stopping it now")
					s.Unlock()
					return
				}

				// send this image to all of the subs
				for c := range s.subs {
					select {
					case c <- jpeg:
					default:
					}
				}
				s.Unlock()

				if avgFrameSize == 0 {
					avgFrameSize = len(jpeg)
				} else {
					avgFrameSize = (avgFrameSize + len(jpeg)) / 2
				}

				frameCount++
			case err, ok := <-errors:
				if !ok {
					return
				}

				log.Warn("unable to get frame", zap.Error(err))

				errCount++
				if errCount >= 24 {
					log.Warn("stopping stream", zap.String("error", "exceeded consecutive error count"))
					return
				}
			}
		}
	}()

	return s, nil
}
