package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/byuoitav/auth/session/cookiestore"
	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// claimAuthorizedRooms is a map[string]map[string]string
	// where first map key is the roomID
	// the second map key is the controlGroup
	// and the value is the time (formatted in RFC3339) they were authenticated
	claimAuthorizedRooms = "rooms"

	// claimAuth is a map[string]bool
	claimAuth = "auth"
)

type ControlHandlers struct {
	ConfigService     cameraservices.ConfigService
	ControlKeyService cameraservices.ControlKeyService
	AuthService       cameraservices.AuthService
	Me                *url.URL
	Logger            *zap.Logger
	SessionStore      *cookiestore.Store
	SessionName       string
	DisableAuth       bool
}

func (h *ControlHandlers) GetCameras(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	var info cameraservices.ControlInfo
	if err := c.BindQuery(&info); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if !h.DisableAuth && !h.authorized(c, info.Room, info.ControlGroup) {
		c.Status(http.StatusUnauthorized)
		return
	}

	cameras, err := h.ConfigService.Cameras(ctx, info)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get cameras: %s", err))
		return
	}

	// change urls to go through proxy (me)
	rewrite := func(u string) string {
		url, err := url.Parse(u)
		if err != nil {
			return ""
		}

		url.Scheme = h.Me.Scheme
		url.Host = h.Me.Host

		switch {
		case strings.Contains(u, "aver"):
			url.Path = "/proxy/aver" + url.Path
		case strings.Contains(u, "axis"):
			url.Path = "/proxy/axis" + url.Path
		}

		return url.String()
	}

	for i := range cameras {
		cameras[i].PanLeft = rewrite(cameras[i].PanLeft)
		cameras[i].PanRight = rewrite(cameras[i].PanRight)
		cameras[i].TiltUp = rewrite(cameras[i].TiltUp)
		cameras[i].TiltDown = rewrite(cameras[i].TiltDown)
		cameras[i].PanTiltStop = rewrite(cameras[i].PanTiltStop)
		cameras[i].ZoomIn = rewrite(cameras[i].ZoomIn)
		cameras[i].ZoomOut = rewrite(cameras[i].ZoomOut)
		cameras[i].ZoomStop = rewrite(cameras[i].ZoomStop)
		cameras[i].Stream = rewrite(cameras[i].Stream)
		cameras[i].Reboot = rewrite(cameras[i].Reboot)

		for j := range cameras[i].Presets {
			cameras[i].Presets[j].SetPreset = rewrite(cameras[i].Presets[j].SetPreset)
		}
	}

	c.JSON(http.StatusOK, cameras)
}

func (h *ControlHandlers) GetControlInfo(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	room, cg, err := h.ControlKeyService.RoomAndControlGroup(ctx, c.Query("key"))
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("unable to get room: %s", err))
		return
	}

	if !h.DisableAuth {
		session, err := h.SessionStore.Get(c.Request, h.SessionName)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			return
		}

		// stick auth map into the cookie
		session.Values[claimAuth] = cameraservices.CtxAuth(ctx)

		if rooms, ok := session.Values[claimAuthorizedRooms].(map[string]interface{}); ok {
			// rooms: map[string]interface{} {"roomID": something}
			if cgs, ok := rooms[room].(map[string]interface{}); ok {
				// cgs: map[string]interface{} {"controlGroup": something}
				if len(cgs) == 0 {
					rooms[room] = map[string]interface{}{
						cg: time.Now().Format(time.RFC3339),
					}
				} else {
					cgs[cg] = time.Now().Format(time.RFC3339)
				}
			} else {
				rooms[room] = map[string]interface{}{
					cg: time.Now().Format(time.RFC3339),
				}
			}

			session.Values[claimAuthorizedRooms] = rooms
		} else {
			session.Values[claimAuthorizedRooms] = map[string]interface{}{
				room: map[string]interface{}{
					cg: time.Now().Format(time.RFC3339),
				},
			}
		}

		_ = session.Save(c.Request, c.Writer)
	}

	c.JSON(http.StatusOK, cameraservices.ControlInfo{
		Room:         room,
		ControlGroup: cg,
	})
}

func (h *ControlHandlers) Proxy(to *url.URL) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO should i validate that they are allowed to get the info for each specific camera?
		id := c.GetString(_cRequestID)

		log := h.Logger
		if len(id) > 0 {
			log = log.With(zap.String("requestID", id))
		}

		defer func() {
			if err := recover(); err != nil {
				if err == http.ErrAbortHandler {
					return
				}

				panic(err)
			}
		}()

		proxy := httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = to.Scheme
				req.URL.Host = to.Host

				switch {
				case strings.HasPrefix(req.URL.Path, "/proxy/aver"):
					req.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/proxy/aver")
				case strings.HasPrefix(req.URL.Path, "/proxy/axis"):
					req.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/proxy/axis")
				}

				req.Header.Set(_hRequestID, id)

				log.Debug("Forwarding request to", zap.String("url", req.URL.String()))
			},
			ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
				log.Warn("error proxying request", zap.Error(err))
				rw.WriteHeader(http.StatusBadGateway)
				_, _ = rw.Write([]byte(fmt.Sprintf("unable to proxy request: %s", err)))
			},
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (h *ControlHandlers) AuthorizeProxy(c *gin.Context) {
	var authorized bool
	path := c.Request.URL.Path

	switch {
	case strings.Contains(path, "reboot"):
		authorized = h.AuthService.IsAuthorizedFor(c.Request.Context(), "reboot")
	case strings.Contains(path, "setPreset"):
		authorized = h.AuthService.IsAuthorizedFor(c.Request.Context(), "setPreset")
	default:
		authorized = h.AuthService.IsAuthorizedFor(c.Request.Context(), "allow")
	}

	if !authorized {
		c.String(http.StatusForbidden, "Unauthorized")
		c.Abort()
		return
	}

	c.Next()
}

func (h *ControlHandlers) authorized(c *gin.Context, room, controlGroup string) bool {
	var authorized bool

	session, err := h.SessionStore.Get(c.Request, h.SessionName)
	if err != nil {
		return authorized
	}

	defer session.Save(c.Request, c.Writer) // nolint:errcheck

	if rooms, ok := session.Values[claimAuthorizedRooms].(map[string]interface{}); ok {
		for roomID, cgs := range rooms {
			if controlGroups, ok := cgs.(map[string]interface{}); ok {
				for cg, v := range controlGroups {
					if ts, ok := v.(string); ok {
						if created, err := time.Parse(time.RFC3339, ts); err == nil {
							switch {
							case time.Since(created).Hours() >= 8:
								delete(controlGroups, cg)
							case room == roomID && controlGroup == cg:
								authorized = true
							}
						} else {
							// not a valid time
							delete(controlGroups, cg)
						}
					} else {
						// not a string
						delete(controlGroups, cg)
					}
				}

				if len(controlGroups) == 0 {
					delete(rooms, roomID)
				}
			}
		}

		if len(rooms) == 0 {
			delete(session.Values, claimAuthorizedRooms)
		}
	}

	return authorized
}
