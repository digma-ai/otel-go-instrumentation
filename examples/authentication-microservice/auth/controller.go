package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/trace"
)

type AuthController struct {
	service Service
	tracer  trace.Tracer
}

func NewAuthController(service Service, tracer trace.Tracer) AuthController {
	return AuthController{
		service: service,
		tracer:  tracer,
	}
}

func (controller *AuthController) Authenticate(c echo.Context) error {
	c.String(http.StatusOK, "Authenticated!")
	return nil
}
