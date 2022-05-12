package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/digma-ai/opentelemetry-go-instrumentation/examples/authentication-microservice/auth"
	"github.com/digma-ai/opentelemetry-go-instrumentation/examples/authentication-microservice/infrastructure/opentelemetry"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
)

var (
	port    = "8011"
	appName = "authentication-microservice"
)

func main() {
	// injected latency
	if s := os.Getenv("EXTRA_LATENCY"); s != "" {
		v, err := time.ParseDuration(s)
		if err != nil {
			log.Fatalf("failed to parse EXTRA_LATENCY (%s) as time.Duration: %+v", v, err) //%+v: variant will include the struct’s field names.
		}
		auth.ExtraLatency = v
		log.Printf("extra latency enabled (duration: %v)", v)
	} else {
		auth.ExtraLatency = time.Duration(0)
	}

	shutdown := opentelemetry.InitTracer()
	defer shutdown()

	tracer := otel.Tracer(appName)

	service := auth.NewAuthService()
	service.Init()
	controller := auth.NewAuthController(service, tracer)

	r := echo.New()
	r.Use(otelecho.Middleware(appName))
	r.POST("/auth", controller.Authenticate)

	fmt.Println("listening on :" + port)
	handleErr(r.Start(":"+port), "failed to listen & serve")
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}
