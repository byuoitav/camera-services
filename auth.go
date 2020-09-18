package cameraservices

import (
	"context"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	FillAuth(*gin.Context)
	AuthorizeFor(...string) gin.HandlerFunc
	IsAuthorizedFor(context.Context, ...string) bool
}

type contextKey string

const (
	authMap contextKey = "authMap"
)

func CtxAuth(ctx context.Context) map[string]bool {
	auth, ok := ctx.Value(authMap).(map[string]bool)
	if !ok {
		return nil
	}

	return auth
}

func WithAuth(ctx context.Context, auth map[string]bool) context.Context {
	return context.WithValue(ctx, authMap, auth)
}
