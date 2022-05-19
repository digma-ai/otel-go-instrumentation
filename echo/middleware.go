package echo

import (
	"errors"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			name, err := getHandlerName(c)
			if err == nil {
				span := trace.SpanFromContext(c.Request().Context())
				span.SetAttributes(attribute.String("endpoint.function_full_name", name))
			}
			return next(c)
		}
	}
}

func getHandlerName(c echo.Context) (string, error) {
	for _, r := range c.Echo().Routes() {
		if r.Method == c.Request().Method && r.Path == c.Path() {
			return r.Name, nil
		}
	}
	return "", errors.New("Handler not found")
}
