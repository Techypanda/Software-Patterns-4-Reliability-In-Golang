package patterns

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/internal/logging"
	"techytechster.com/softwarepatterns/model"
)

func isSuccess(state *model.State) bool {
	roll := rand.Intn(100-0+1) + 0
	return roll > int(state.FailureChance) && state.FailureChance < 100
}

func Failure(c echo.Context) error {
	logging.WriteLogLn("Failed To Succeed - Returning Sad Image")
	path, _ := filepath.Abs("../../data/failure")
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println("failure to read img for real", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failure to read img for real",
		})
	}
	pic := files[rand.Intn(len(files))]
	return c.File(fmt.Sprintf("%s/%s", path, pic.Name()))
}

func PictureAPI(c echo.Context) error {
	logging.WriteLogLn("Succeeded! - Returning Success Image")
	path, _ := filepath.Abs("../../data/successimgs")
	files, err := os.ReadDir(path)
	if err != nil {
		log.Println("failure to read img for real", err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "failure to read img for real",
		})
	}
	pic := files[rand.Intn(len(files))]
	return c.File(fmt.Sprintf("%s/%s", path, pic.Name()))
}
