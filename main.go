package main

import (
	"armiariyan/attendances-system-jwt/config"
	"armiariyan/attendances-system-jwt/controller"
	"armiariyan/attendances-system-jwt/helper"
	"armiariyan/attendances-system-jwt/middleware"
	"armiariyan/attendances-system-jwt/repository"
	"armiariyan/attendances-system-jwt/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                  = config.SetupDatabaseConnection()
	userRepository repository.UserRepository = repository.NewUserRepository(db)
	jwtService     service.JWTService        = service.NewJWTService()
	userService    service.UserService       = service.NewUserService(userRepository)
	userController controller.UserController = controller.NewUserController(userService, jwtService)

	authMiddleware    gin.HandlerFunc = middleware.AuthorizeJWT(userService, jwtService)
	checkInMiddleware gin.HandlerFunc = middleware.CheckInMiddleware(userService)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	// seeder.DBSeed(db) //<-- this will seed 2 user, check seeder.go for the credential

	authRoutes := r.Group("api/v1")
	{
		authRoutes.POST("/register", userController.Register)
		authRoutes.POST("/login", userController.Login)
	}

	attendanceRoutes := r.Group("api/v1", authMiddleware)
	{
		attendanceRoutes.POST("/checkin", userController.CheckIn)
		attendanceRoutes.POST("/checkout", userController.CheckOut)
		attendanceRoutes.GET("/attendances", userController.GetAttendancesHistory)
	}

	activityRoutes := r.Group("api/v1", authMiddleware, checkInMiddleware)
	{
		activityRoutes.POST("/activity", checkInMiddleware, userController.CreateActivity)
		activityRoutes.PUT("/activity/:id_activity", checkInMiddleware, userController.UpdateActivity)
		activityRoutes.DELETE("/activity/:id_activity", checkInMiddleware, userController.DeleteActivity)
		activityRoutes.GET("/activity/:id", userController.GetActivityHistoryByDate)
	}

	checkRoutes := r.Group("api/v1")
	{
		checkRoutes.GET("/health", func(c *gin.Context) {
			response := helper.BuildResponse(true, "healthcheck successfull", nil)
			c.JSON(http.StatusOK, response)
		})
	}

	r.Run()
}
