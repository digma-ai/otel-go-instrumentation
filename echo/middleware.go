package echo

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Middleware(e *echo.Echo) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//key := c.Request().Method+c.Request().RequestURI
			for _, r := range e.Routes() {
				if r.Method == c.Request().Method && r.Path == c.Path() {
					fmt.Println(r.Name)
				}
			}
			return next(c)
		}
	}
}
