package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/internal/logging"
	"techytechster.com/softwarepatterns/internal/patterns"
	"techytechster.com/softwarepatterns/model"
)

type SwitchPatternPayload struct {
	Mode string `json:"pattern" form:"pattern" query:"pattern" validate:"required"`
}
type UpdateFailureChancePayload struct {
	Chance float64 `json:"chance" form:"chance" query:"chance" validate:"required,min=-1,max=100"` // Technically you can go -0.01, but it doesnt break the program we consider <=0 == 0
}

func ReadLogs(c echo.Context) error {
	logContents, err := logging.ReadLogContents()
	if err != nil {
		log.Println("failed to read file! - ", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "failed to read log",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"contents": logContents,
	})
}

func DyingAPICall(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		var toCall func(c echo.Context) error
		switch state.Mode {
		case model.BACKOFF:
			toCall = patterns.Backoff(state)
		case model.EXPONENTIAL_BACKOFF:
			toCall = patterns.ExponentialBackoff(state)
		case model.EXP_WITH_JITTER:
			toCall = patterns.ExponentialBackoffWithJitter(state)
		case model.CIRCUIT_BREAKER:
			toCall = patterns.DoCircuitBreaking(state)
		default:
			toCall = patterns.NaiivePattern(state)
		}
		return toCall(c)
	}
}

func ChangeFailureChance(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(UpdateFailureChancePayload)
		if err := c.Bind(payload); err != nil {
			log.Println("failed to write incoming request to payload struct")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "failed to write incoming request to payload struct",
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("failed to validate to payload: %s\n", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("Expected failureChance: %s", err.Error()),
			})
		}
		state.FailureChance = payload.Chance
		logging.WriteLogLn(fmt.Sprintf("Updated chance to %f", payload.Chance))
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"failureChance": payload.Chance,
		})
	}
}

func CheckFailureChance(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"failureChance": state.FailureChance,
		})
	}
}

func SwitchPattern(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		payload := new(SwitchPatternPayload)
		if err := c.Bind(payload); err != nil {
			log.Println("failed to write incoming request to payload struct")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "failed to write incoming request to payload struct",
			})
		} else if err := c.Validate(payload); err != nil {
			log.Printf("failed to validate to payload: %s\n", err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("Expected state: %s", err.Error()),
			})
		}
		if payload.Mode != model.DEFAULT_MODE && payload.Mode != model.BACKOFF && payload.Mode != model.EXPONENTIAL_BACKOFF && payload.Mode != model.EXP_WITH_JITTER && payload.Mode != model.ALL && payload.Mode != model.CIRCUIT_BREAKER {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Expected NO_PATTERN or EXPONENTIAL_BACKOFF or EXP_WITH_JITTER or ALL or CIRCUIT_BREAKER",
			})
		}
		state.Mode = payload.Mode
		logging.WriteLogLn(fmt.Sprintf("Successfully Updated Pattern To %s", payload.Mode))
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"success": "Successfully Updated State",
		})
	}
}

func CurrentPattern(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"currentPattern": state.Mode,
		})
	}
}

func Heartbeat(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"ok": true,
	})
}
