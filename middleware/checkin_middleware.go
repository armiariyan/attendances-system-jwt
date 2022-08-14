package middleware

import (
	"armiariyan/attendances-system-jwt/entity"
	"armiariyan/attendances-system-jwt/helper"
	"armiariyan/attendances-system-jwt/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

//CheckInMiddleware
func CheckInMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(entity.User)

		// Cek if user already check in today
		ok, err := userService.IsCheckInToday(currentUser.Id)
		if !ok || err != nil {
			// Build response error
			response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}
	}
}
