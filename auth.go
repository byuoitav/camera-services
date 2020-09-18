package cameraservices

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	FillAuth(*gin.Context)
	AuthorizeFor(...string) gin.HandlerFunc
	IsAuthorizedFor(context.Context, ...string) bool
	Auth(context.Context) map[string]bool
}
