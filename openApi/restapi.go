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
)

const serverAddress = "localhost:8080"

type SeptimanappRestApi struct {
	repository database.Repository
}

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

func (s SeptimanappRestApi) AuthorizeModifyEvents(key string, _ echo.Context) (bool, error) {
	return security.ValidateApikey(s.repository, key)
}

func (s SeptimanappRestApi) GetEvents(ctx echo.Context, params GetEventsParams) error {
	events, err := s.repository.GetEvents(params.Year)
	if err == nil {
		return ctx.JSON(http.StatusOK, events)
	} else {
		return sendInternalError(ctx)
	}
}

func (s SeptimanappRestApi) GetEventsId(ctx echo.Context, id EventId) error {
	event, err := s.repository.GetEvent(int(id))
	if err == nil {
		return ctx.JSON(http.StatusOK, event)
	} else {
		return sendInternalError(ctx)
	}
}

func (s SeptimanappRestApi) PostEvents(ctx echo.Context) error {
	var events types.Events
	err := ctx.Bind(&events)
	if err != nil {
		return sendInvalidRequest(ctx)
	}
	err = ctx.Validate(events)
	if err != nil {
		fmt.Println("Not validated")
		fmt.Println(err)
		return sendInvalidFormat(ctx, "events")
	}
	ids, err := s.repository.AddEvents(events)
	if err != nil {
		return sendInternalError(ctx)
	}
	fmt.Println("POSTED:")
	fmt.Println(ids)
	return ctx.JSON(http.StatusOK, ids)
}

func sendInvalidFormat(ctx echo.Context, paramName string) error {
	return ctx.String(http.StatusBadRequest, fmt.Sprintf("Invalid format for %s", paramName))
}

func sendInvalidRequest(ctx echo.Context) error {
	return ctx.String(http.StatusBadRequest, "Invalid Request")
}

func (s SeptimanappRestApi) DeleteEventsId(ctx echo.Context, id EventId) error {
	err := s.repository.DeleteEvent(int(id))
	if err != nil {
		return sendInternalError(ctx)
	}
	return sendOK(ctx)
}

func (s SeptimanappRestApi) PutEventsId(ctx echo.Context, id EventId) error {
	var event types.Event
	err := ctx.Bind(&event)
	if err != nil {
		return sendInvalidRequest(ctx)
	}
	err = ctx.Validate(event)
	if err != nil {
		fmt.Println("Not validated")
		fmt.Println(err)
		return sendInvalidFormat(ctx, "event")
	}
	event.ID = int(id)
	err = s.repository.UpdateEvent(event)
	if err != nil {
		if err.Error() == database.RecordNotFound {
			return ctx.String(http.StatusNotFound, "No such event")
		}
		return sendInternalError(ctx)
	}
	fmt.Printf("UPDATED:%v\n", id)
	return sendOK(ctx)
}

func (s SeptimanappRestApi) GetLocations(ctx echo.Context, params GetLocationsParams) error {
	location, err := s.repository.GetLocations(params.OverallLocation)
	if err == nil {
		return ctx.JSON(http.StatusOK, location)
	} else {
		return sendInternalError(ctx)
	}
}

func (s SeptimanappRestApi) GetLocationsId(ctx echo.Context, id string) error {
	location, err := s.repository.GetLocation(id)
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
	repository, err := database.GetRepository()
	if err != nil {
		panic(err)
	}
	api := SeptimanappRestApi{repository: repository}
	RegisterHandlers(e, api)
	RegisterAuthMiddleware(e, api) // TODO remove if codegen provides the corresponding part
}

func StartRestApi() {
	fmt.Println("START REST:")
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}

	//e.Use(echomiddleware.KeyAuthWithConfig(echomiddleware.KeyAuthConfig{
	//	KeyLookup: "query:appid",
	//	Skipper: func(ctx echo.Context) bool {
	//		fmt.Println("keyauth:")
	//		fmt.Println(ctx.Get(App_idScopes))
	//		if strings.HasPrefix(ctx.Path(), "/openapi/definition") || strings.HasPrefix(ctx.Path(), "/definition") {
	//			return true
	//		}
	//		return false
	//	},
	//	Validator:  func(key string, ctx echo.Context) (bool, error) {
	//		return key == "valid-key", nil
	//	},
	//}))
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
