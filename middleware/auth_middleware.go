package middleware

import (
	"armiariyan/attendances-system-jwt/helper"
	"armiariyan/attendances-system-jwt/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

//AuthorizeJWT validates the token user given, return 401 if not valid
func AuthorizeJWT(userService service.UserService, jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("failed to process request", "no token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Validate token
		token, err := jwtService.ValidateToken(authHeader)
		// If token not valid
		if !token.Valid && err != nil {
			log.Println(err)
			response := helper.BuildErrorResponse("token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Take all claims/payload from jwt
		claims := token.Claims.(jwt.MapClaims)

		// Check if token expires
		expiresAt := claims["standard_claims"].(map[string]interface{})["exp"].(float64)
		if int64(expiresAt) < time.Now().Unix() {
			response := helper.BuildErrorResponse("token is not valid", "token expired!", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		// Get User By ID, claims type is interface
		user_id, _ := strconv.Atoi(claims["user_id"].(string))
		user, _ := userService.GetUserById(user_id)

		// Set context for current user
		c.Set("currentUser", user)
	}
}
