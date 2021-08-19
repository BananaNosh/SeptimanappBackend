package Openapi

import "C"
import (
	"SeptimanappBackend/database"
	"SeptimanappBackend/security"
	"SeptimanappBackend/types"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"net/http"
	"reflect"
	"strings"
	"time"
)

const serverAddress = "localhost:8080"

type SeptimanappRestApi struct{}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if reflect.TypeOf(i).Kind() == reflect.Slice {
		if reflect.ValueOf(i).Len() > 0 && reflect.ValueOf(i).Index(0).FieldByName("ID").IsValid() {
			return v.validator.Var(i, "required,unique=ID,dive")
		}
		return v.validator.Var(i, "required,min=1,dive")
		//s := reflect.ValueOf(i)
		//for i := 0; i < s.Len(); i++ {
		//	if err := v.validator.Struct(s.Index(i).Interface()); err != nil {
		//		return err
		//	}
		//}
		//return nil
	} else {
		return v.validator.Struct(i)
	}
}

func sendOK(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "OK")
}

func sendInternalError(ctx echo.Context) error {
	return ctx.String(http.StatusInternalServerError, "There was an error with the database")
}

func (s SeptimanappRestApi) GetEvents(ctx echo.Context, params GetEventsParams) error {
	events, err := database.GetEvents(params.Year)
	if err == nil {
		return ctx.JSON(http.StatusOK, events)
	} else {
		return sendInternalError(ctx)
	}
}

func (s SeptimanappRestApi) PostEvents(ctx echo.Context) error {
	var events []types.Event
	err := ctx.Bind(&events)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid format for events")
	}
	err = ctx.Validate(events)
	fmt.Println("POSTED:")
	fmt.Println(events)
	if err != nil {
		fmt.Println("Not validated")
		fmt.Println(err)
		return ctx.String(http.StatusBadRequest, "Invalid format for events")
	}
	return sendOK(ctx)
}

func (s SeptimanappRestApi) AuthorizePostEvents(key string, ctx echo.Context) (bool, error) {
	return security.ValidateApikey(key)
}

func (s SeptimanappRestApi) GetEventsId(ctx echo.Context, id int) error {
	event, err := database.GetEvent(id)
	if err == nil {
		return ctx.JSON(http.StatusOK, event)
	} else {
		return sendInternalError(ctx)
	}
}

func (s SeptimanappRestApi) GetLocations(ctx echo.Context, params GetLocationsParams) error {
	location, err := database.GetLocations(params.OverallLocation)
	if err == nil {
		return ctx.JSON(http.StatusOK, location)
	} else {
		return sendInternalError(ctx)
	}
}

func (s SeptimanappRestApi) GetLocationsId(ctx echo.Context, id string) error {
	location, err := database.GetLocation(id)
	if err == nil {
		return ctx.JSON(http.StatusOK, location)
	} else {
		return sendInternalError(ctx)
	}
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
	api := SeptimanappRestApi{}
	RegisterHandlers(e, api)
	RegisterAuthMiddleware(e, api) // TODO remove if codegen provides the corresponding part
}

func StartRestApi() {
	start := time.Now()
	year := 2020
	events, err := database.GetEvents(&year)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print("Events\n")
		fmt.Println(len(events))
		if len(events) > 0 {
			fmt.Println(events[0].Start)
		}
	}

	fmt.Printf("%d mikro s\n", time.Since(start)/1000)
	fmt.Println("START REST:")
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	swagger, err := GetSwagger()
	if err != nil {
		panic(err)
	}
	e.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Skipper: func(ctx echo.Context) bool {
			//print(ctx.Path())
			if strings.HasPrefix(ctx.Path(), "/openapi/definition") || strings.HasPrefix(ctx.Path(), "/definition") {
				return true
			}
			//return strings.HasPrefix(ctx.Request().Host, "localhost")
			return false
		},
		Options: openapi3filter.Options{
			AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
		},
	}))

	SetupDocumentationRoutes(e)
	SetupRestRoutes(e)

	e.Logger.Fatal(e.Start(serverAddress))
}
