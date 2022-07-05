package patterns

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/internal/logging"
	"techytechster.com/softwarepatterns/model"
)

func ExponentialBackoffWithJitter(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		retryCount := 5
		timeMultiplier := 4.0
		for x := 0; x < retryCount; x++ {
			if isSuccess(state) {
				return PictureAPI(c)
			}
			timeMultiplier *= timeMultiplier              // ^2
			minVal := int(math.Min(5000, timeMultiplier)) // 5 Seconds AT MOST
			jittered := rand.Intn(minVal-0+1) + 0         // random between 0 and min(MAX_WE_WANT_TO_WAIT, Exponential)
			logging.WriteLogLn(fmt.Sprintf("Failed To Get A Success, Using Backoff Going To Retry In %dms, %d tries remaining", jittered, retryCount-(x+1)))
			time.Sleep(time.Duration(jittered) * time.Millisecond)
		}
		logging.WriteLogLn("Failed To Get A Success After Retrying 5 Times, Giving Up")
		return Failure(c)
	}
}
