package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/byuoitav/auth/wso2"
	"github.com/byuoitav/camera-services/cmd/control/couch"
	"github.com/byuoitav/camera-services/cmd/control/keys"
	"github.com/byuoitav/camera-services/handlers"
	"github.com/byuoitav/common/v2/auth"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var (
		port     int
		logLevel string

		dbAddr     string
		dbUsername string
		dbPassword string
		dbInsecure bool

		keyServiceAddr string
	)

	pflag.CommandLine.IntVarP(&port, "port", "P", 8080, "port to run the server on")
	pflag.StringVarP(&logLevel, "log-level", "L", "", "level to log at. refer to https://godoc.org/go.uber.org/zap/zapcore#Level for options")
	pflag.StringVar(&dbAddr, "db-address", "", "database address")
	pflag.StringVar(&dbUsername, "db-username", "", "database username")
	pflag.StringVar(&dbPassword, "db-password", "", "database password")
	pflag.BoolVar(&dbInsecure, "db-insecure", false, "don't use SSL in database connection")
	pflag.StringVar(&keyServiceAddr, "key-service", "control-keys.av.byu.edu", "address of the control keys service")
	pflag.Parse()

	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		fmt.Printf("invalid log level: %s\n", err.Error())
		os.Exit(1)
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json", EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	log, err := config.Build()
	if err != nil {
		fmt.Printf("unable to build logger: %s", err)
		os.Exit(1)
	}
	defer func() {
		_ = log.Sync()
	}()

	// context for setup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// build the config service
	if dbInsecure {
		dbAddr = "http://" + dbAddr
	} else {
		dbAddr = "https://" + dbAddr
	}

	client := wso2.Client{
		CallbackURL:  os.Getenv("CALLBACK_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		GatewayURL:   os.Getenv("GATEWAY_URL"),
	}
	writeconfig := router.Group(
		"",
		auth.CheckHeaderBasedAuth,
		gin.WrapMiddleware(client.AuthCodeMiddleware),
		auth.AuthorizeRequest("write-config", "configuration", func(c echo.Context) string { return "all" }),
	)
	readconfig := router.Group(
		"",
		auth.CheckHeaderBasedAuth,
		gin.WrapMiddleware(client.AuthCodeMiddleware),
		auth.AuthorizeRequest("read-config", "configuration", func(c echo.Context) string { return "all" }),
	)

	var csOpts []couch.Option
	if dbUsername != "" {
		csOpts = append(csOpts, couch.WithBasicAuth(dbUsername, dbPassword))
	}
	fmt.Printf("db addr: %s\n", dbAddr)
	cs, err := couch.New(ctx, dbAddr, csOpts...)
	if err != nil {
		log.Fatal("unable to create config service", zap.Error(err))
	}

	handlers := handlers.ControlHandlers{
		ConfigService: cs,
		ControlKeyService: &keys.ControlKeyService{
			Address: keyServiceAddr,
		},
	}

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/key/:key", handlers.GetCameras)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("unable to bind listener", zap.Error(err))
	}

	log.Info("Starting server", zap.String("on", lis.Addr().String()))
	err = r.RunListener(lis)
	switch {
	case errors.Is(err, http.ErrServerClosed):
	case err != nil:
		log.Fatal("failed to server", zap.Error(err))
	}
}