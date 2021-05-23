package Openapi

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

const serverAddress = "localhost:8080"

func StartRestApi() {
	e := echo.New()
	//e.Pre(middleware.Rewrite(map[string]string{ TODO check
	//	//"/openapi/definition": "/openapi/definition/index.html",
	//	"/openapi/definition/test": "/openapi/definition/index.html",
	//}))
	e.Static("/definition", "./openApi/definition")
	e.GET("/openapi/definition/*", echoSwagger.EchoWrapHandler(echoSwagger.URL(fmt.Sprintf("http://%s/definition/openapi.yaml", serverAddress))))
	e.Logger.Fatal(e.Start(serverAddress))
}
