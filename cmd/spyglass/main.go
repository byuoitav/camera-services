package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/byuoitav/auth/session/cookiestore"
	"github.com/byuoitav/auth/wso2"
	"github.com/byuoitav/camera-services/couch"
	"github.com/byuoitav/camera-services/keys"
	"github.com/byuoitav/camera-services/opa"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

const (
	sessionName = "camera-services-spyglass"
)

func main() {
	var (
		port int

		dbAddr     string
		dbUsername string
		dbPassword string
		dbInsecure bool

		keyServiceAddr string

		callbackURL  string
		clientID     string
		clientSecret string
		gatewayURL   string

		opaURL      string
		opaToken    string
		disableAuth bool

		controlURLFormat string
	)

	pflag.CommandLine.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVar(&dbAddr, "db-address", "", "database address")
	pflag.StringVar(&dbUsername, "db-username", "", "database username")
	pflag.StringVar(&dbPassword, "db-password", "", "database password")
	pflag.BoolVar(&dbInsecure, "db-insecure", false, "don't use SSL in database connection")
	pflag.StringVar(&keyServiceAddr, "key-service", "control-keys.av.byu.edu", "address of the control keys service")
	pflag.StringVar(&callbackURL, "callback-url", "http://localhost:8080", "wso2 callback url")
	pflag.StringVar(&clientID, "client-id", "", "wso2 client ID")
	pflag.StringVar(&clientSecret, "client-secret", "", "wso2 client secret")
	pflag.StringVar(&gatewayURL, "gateway-url", "https://api.byu.edu", "ws02 gateway url")
	pflag.StringVar(&opaURL, "opa-url", "", "The URL of the OPA Authorization server")
	pflag.StringVar(&opaToken, "opa-token", "", "The token to use for OPA")
	pflag.BoolVar(&disableAuth, "disable-auth", false, "Disable all auth z/n checks")
	pflag.StringVar(&controlURLFormat, "control-url", "https://cameras.av.byu.edu/login?key=%s", "The url format string of the camera control service")

	pflag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build the config service
	if dbInsecure {
		dbAddr = "http://" + dbAddr
	} else {
		dbAddr = "https://" + dbAddr
	}

	var csOpts []couch.Option
	if dbUsername != "" {
		csOpts = append(csOpts, couch.WithBasicAuth(dbUsername, dbPassword))
	}

	cs, err := couch.New(ctx, dbAddr, csOpts...)
	if err != nil {
		log.Fatalf("unable to create config service: %s", err)
	}

	handlers := Handlers{
		CameraControlURLFormat: controlURLFormat,
		ConfigService:          cs,
		ControlKeyService: &keys.ControlKeyService{
			Address: keyServiceAddr,
		},
	}

	wso2 := wso2.New(clientID, clientSecret, gatewayURL, callbackURL)
	sessionStore := cookiestore.NewStore()

	auth := opa.Client{
		Address:  opaURL,
		Endpoint: "/v1/data/spyglass",
		Token:    opaToken,
	}

	r := gin.New()
	r.Use(cors.Default())
	r.Use(gin.Recovery())

	if !disableAuth {
		r.Use(adapter.Wrap(wso2.AuthCodeMiddleware(sessionStore, sessionName)))
		r.Use(auth.Authorize)
	}

	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)

		if file == "" || filepath.Ext(file) == "" {
			c.File("/web/index.html")
		} else {
			c.File("/web/" + path.Join(dir, file))
		}
	})

	api := r.Group("/api/v1/")
	api.GET("/rooms", handlers.GetRooms)
	api.GET("/rooms/:room/controlGroups", handlers.GetControlGroups)
	api.GET("/rooms/:room/controlGroups/:controlGroup", handlers.ControlPage)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("unable to bind listener", zap.Error(err))
	}

	log.Printf("Starting server on %s", lis.Addr().String())
	err = r.RunListener(lis)
	switch {
	case errors.Is(err, http.ErrServerClosed):
	case err != nil:
		log.Fatal("failed to server", zap.Error(err))
	}
}
