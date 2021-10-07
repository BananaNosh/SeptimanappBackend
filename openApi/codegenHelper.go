package Openapi

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func (s *SeptimanappRestApi) ProvideAuthMiddlewareModifyEvents() echo.MiddlewareFunc {
	return echomiddleware.KeyAuthWithConfig(echomiddleware.KeyAuthConfig{
		KeyLookup: "query:appid",
		Validator: s.AuthorizeModifyEvents,
		Skipper: func(ctx echo.Context) bool {
			return ctx.Request().Method == "GET"
		},
	})
}

func RegisterAuthMiddleware(e *echo.Echo, s SeptimanappRestApi) {
	e.Use(s.ProvideAuthMiddlewareModifyEvents())
}
