package Openapi

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

const serverAddress = "localhost:8080"

type SeptimanappRestApi struct{}

func (s SeptimanappRestApi) GetEvents(ctx echo.Context, params GetEventsParams) error {
	panic("implement me")
}

func (s SeptimanappRestApi) PostEvents(ctx echo.Context) error {
	panic("implement me")
}

func (s SeptimanappRestApi) GetEventsId(ctx echo.Context, id int) error {
	panic("implement me")
}

func (s SeptimanappRestApi) GetLocations(ctx echo.Context) error {
	panic("implement me")
}

func (s SeptimanappRestApi) GetLocationsId(ctx echo.Context, id string) error {
	panic("implement me")
}

func SetupDocumentationRoutes(e *echo.Echo) {
	//e.Pre(middleware.Rewrite(map[string]string{ TODO check
	//	//"/openapi/definition": "/openapi/definition/index.html",
	//	"/openapi/definition/test": "/openapi/definition/index.html",
	//}))
	e.Static("/definition", "./openApi/definition")
	e.GET("/openapi/definition/*", echoSwagger.EchoWrapHandler(echoSwagger.URL(fmt.Sprintf("http://%s/definition/openapi.yaml", serverAddress))))
}

func SetupRestRoutes(e *echo.Echo) {
	RegisterHandlers(e, SeptimanappRestApi{})
}

func StartRestApi() {
	e := echo.New()
	SetupDocumentationRoutes(e)
	e.Logger.Fatal(e.Start(serverAddress))
}
