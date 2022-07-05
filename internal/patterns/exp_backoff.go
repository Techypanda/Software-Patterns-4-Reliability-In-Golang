package patterns

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/internal/logging"
	"techytechster.com/softwarepatterns/model"
)

func ExponentialBackoff(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		retryCount := 2
		timeMultiplier := 4.0
		for x := 0; x < retryCount; x++ {
			if isSuccess(state) {
				return PictureAPI(c)
			}
			timeMultiplier *= timeMultiplier // ^2
			logging.WriteLogLn(fmt.Sprintf("Failed To Get A Success, Using Backoff Going To Retry In %fms, %d tries remaining", timeMultiplier, retryCount-(x+1)))
			time.Sleep(time.Duration(timeMultiplier) * time.Millisecond)
		}
		logging.WriteLogLn("Failed To Get A Success After Retrying 5 Times, Giving Up")
		return Failure(c)
	}
}
