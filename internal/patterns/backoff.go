package patterns

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/internal/logging"
	"techytechster.com/softwarepatterns/model"
)

func Backoff(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		fmt.Println("CALLED")
		retryCount := 5
		for x := 0; x < retryCount; x++ {
			if isSuccess(state) {
				return PictureAPI(c)
			}
			logging.WriteLogLn(fmt.Sprintf("Failed To Get A Success, Using Backoff Going To Retry In 300ms, %d tries remaining", retryCount-(x+1)))
			time.Sleep(time.Millisecond * 300)
		}
		logging.WriteLogLn("Failed To Get A Success After Retrying 5 Times, Giving Up")
		return Failure(c)
	}
}
