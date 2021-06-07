package Openapi

import (
	"SeptimanappBackend/database"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"net/http"
	"time"
)

const serverAddress = "localhost:8080"

type SeptimanappRestApi struct{}

func (s SeptimanappRestApi) GetEvents(ctx echo.Context, params GetEventsParams) error {
	events, err := database.GetEvents(params.Year)
	if err == nil {
		return ctx.JSON(http.StatusOK, events)
	} else {
		return ctx.String(500, "There was an error with the database")
	}
}

func (s SeptimanappRestApi) PostEvents(ctx echo.Context) error {
	panic("implement me")
}

func (s SeptimanappRestApi) GetEventsId(ctx echo.Context, id int) error {
	event, err := database.GetEvent(id)
	if err == nil {
		return ctx.JSON(http.StatusOK, event)
	} else {
		return ctx.String(500, "There was an error with the database")
	}
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
	SetupRestRoutes(e)
	e.Logger.Fatal(e.Start(serverAddress))
}
