package echo

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, r := range c.Echo().Routes() {
				if r.Method == c.Request().Method && r.Path == c.Path() {
					span := trace.SpanFromContext(c.Request().Context())
					span.SetAttributes(attribute.String("endpoint.function_full_name", r.Name))
					break
				}
			}
			return next(c)
		}
	}
}
