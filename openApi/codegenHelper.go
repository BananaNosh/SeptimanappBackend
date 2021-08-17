package Openapi

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"strings"
)

func (s *SeptimanappRestApi) ProvideAuthMiddlewarePostEvents() echo.MiddlewareFunc {
	return echomiddleware.KeyAuthWithConfig(echomiddleware.KeyAuthConfig{
		KeyLookup: "query:appid",
		Validator: s.AuthorizePostEvents,
		Skipper: func(ctx echo.Context) bool {
			return !(strings.HasPrefix(ctx.Path(), "/events") && ctx.Request().Method == "POST")
		},
	})
}

func RegisterAuthMiddleware(e *echo.Echo, s SeptimanappRestApi) {
	e.Use(s.ProvideAuthMiddlewarePostEvents())
}
