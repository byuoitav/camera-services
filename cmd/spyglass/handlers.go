package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	CameraControlURLFormat string

	ConfigService interface {
		Rooms(context.Context) ([]string, error)
		ControlGroups(context.Context, string) ([]string, error)
	}

	ControlKeyService interface {
		ControlKey(context.Context, string, string) (string, error)
	}
}

func (h *Handlers) GetRooms(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rooms, err := h.ConfigService.Rooms(ctx)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (h *Handlers) GetControlGroups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	groups, err := h.ConfigService.ControlGroups(ctx, c.Param("room"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (h *Handlers) ControlPage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	key, err := h.ControlKeyService.ControlKey(ctx, c.Param("room"), c.Param("controlGroup"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf(h.CameraControlURLFormat, key))
}
