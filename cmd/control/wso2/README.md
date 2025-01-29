# WSO2 Auth Package

This package provides several different ways of interacting with WSO2, such as
handling OAuth2 when calling APIs behind WSO2, as well as handling the
Authorization Code grant types.

### Middleware

#### Authorization Code Middleware

This middleware allows you to force users to login when hitting an http endpoint
by utilizing the OAuth2 Authorization Code grant type. 

``` golang
package main

import (
	"net/http"

	"github.com/byuoitav/auth/wso2"
	"github.com/byuoitav/auth/session/cookiestore"
	"github.com/gin-gonic/gin"
)

func main() {

	// Create WSO2 Client
	client := wso2.New("client_id", "client_secret", "http://gateway:80", "http://localhost:8080")

	// Create Session Store
	sessionStore := cookiestore.NewStore()

	// Create Gin Router
	router := gin.Default()

	// Utilize Auth Code Middleware
	router.Use(func(c *gin.Context) {
		client.AuthCodeMiddleware(sessionStore, "default-session")(c.Writer, c.Request)
		c.Next()
	})

	router.GET("/hello_world", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	router.Run(":8080")
}
```
