package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"techytechster.com/softwarepatterns/internal"
	"techytechster.com/softwarepatterns/model"
)

type SimpleValidator struct {
	validator *validator.Validate
}

func (cv *SimpleValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	e := echo.New()
	e.Validator = &SimpleValidator{validator: validator.New()}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // horrible practice but this is for a demo DONT DO THIS
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))
	state := model.NewState()
	e.GET("/", internal.Heartbeat)
	e.GET("/log", internal.ReadLogs)
	e.GET("/imageAPI", internal.DyingAPICall(&state))
	e.GET("/chance", internal.CheckFailureChance(&state))
	e.POST("/chance", internal.ChangeFailureChance(&state))
	e.POST("/pattern", internal.SwitchPattern(&state))
	e.GET("/pattern", internal.CurrentPattern(&state))
	e.Logger.Fatal((e.Start(":8080")))
}
