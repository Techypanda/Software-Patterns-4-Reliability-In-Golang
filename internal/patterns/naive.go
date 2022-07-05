package patterns

import (
	"github.com/labstack/echo/v4"
	"techytechster.com/softwarepatterns/model"
)

func NaiivePattern(state *model.State) func(c echo.Context) error {
	return func(c echo.Context) error {
		if isSuccess(state) {
			return PictureAPI(c)
		}
		return Failure(c)
	}
}
