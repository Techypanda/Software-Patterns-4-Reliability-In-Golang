package patterns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/internal/logging"
	"techytechster.com/softwarepatterns/model"
)

type CircuitBreaker struct {
	TTL int64
}

func circuitBreaking(c echo.Context) error {
	img, _ := filepath.Abs("../../data/robot/circuitbreaker.jfif")
	return c.File(img)
}

func verifyIfCircuitBreaker() bool { // Hacky But I Don't Have Time To Do Correctly :(
	jsonFile, _ := filepath.Abs("../../data/robot/break.json")
	contents, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("error occured reading break.json, assuming its empty but going to delete it incase")
		os.Remove(jsonFile)
	}
	breaker := CircuitBreaker{}
	if err := json.Unmarshal(contents, &breaker); err != nil {
		fmt.Println("error occured unmarsheling break.json, assuming its corrupted going to delete it")
		os.Remove(jsonFile)
	}
	if breaker.TTL > time.Now().Unix() {
		logging.WriteLogLn("breaker exists and hasn't expired yet, do robot things")
		return true
	} else {
		fmt.Println("breaker TTL ended, kill it")
		os.Remove(jsonFile)
		return false
	}
}

func NewCircuitBreaker(c echo.Context) error {
	newCircuitBreaker := CircuitBreaker{
		TTL: time.Now().Unix() + 30,
	}
	jsonMarsheled, err := json.Marshal(newCircuitBreaker)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to marshal circuit breaker",
		})
	}
	jsonFile, _ := filepath.Abs("../../data/robot/break.json")
	f, err := os.Create(jsonFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to open circuit breaker",
		})
	}
	_, err = f.WriteString(string(jsonMarsheled))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to save circuit breaker",
		})
	}
	return circuitBreaking(c)
}

func DoCircuitBreaking(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		//return c.String(200, "AAA")
		if verifyIfCircuitBreaker() { // if circuit breaker is still going, use it
			return circuitBreaking(c)
		}
		if isSuccess(state) { // lets check the API again
			return PictureAPI(c)
		}
		return NewCircuitBreaker(c) // its broken lets use a circuit breaker and put it in place
	}
}
