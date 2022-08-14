package controller

import (
	"armiariyan/attendances-system-jwt/dto"
	"armiariyan/attendances-system-jwt/entity"
	"armiariyan/attendances-system-jwt/helper"
	"armiariyan/attendances-system-jwt/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
	// Logout(context *gin.Context)
	CheckIn(context *gin.Context)
	CheckOut(context *gin.Context)
	CreateActivity(context *gin.Context)
	UpdateActivity(context *gin.Context)
	DeleteActivity(context *gin.Context)
	GetActivityHistoryByDate(context *gin.Context)
	GetAttendancesHistory(context *gin.Context)
}

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

func NewUserController(user service.UserService, jwtService service.JWTService) UserController {
	return &userController{
		userService: user,
		jwtService:  jwtService,
	}
}

func (c *userController) Register(context *gin.Context) {
	// Fill RegisterDTO variable for validation
	var registerDTO dto.RegisterDTO
	errDTO := context.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("register failed", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Create User
	createdUser, err := c.userService.CreateUser(registerDTO)
	if err != nil {
		response := helper.BuildErrorResponse("register failed", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Generate JWT token
	token := c.jwtService.GenerateToken(strconv.FormatInt(int64(createdUser.Id), 10))
	createdUser.Token = token

	//Build Response
	response := helper.BuildResponse(true, "OK!", createdUser)
	context.JSON(http.StatusCreated, response)
}

func (c *userController) Login(context *gin.Context) {
	// Fill loginDTO variable
	var loginDTO dto.LoginDTO
	errDTO := context.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Verify the data
	user, err := c.userService.VerifyCredential(loginDTO)
	if err != nil {
		// Build response error
		response := helper.BuildErrorResponse("failed to process request", "invalid email or password", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	// Because the id is int and generate token param is string, so convert it
	user.Token = c.jwtService.GenerateToken(strconv.Itoa(user.Id))

	//Build response if success
	response := helper.BuildResponse(true, "successfully logged in", user)
	context.JSON(http.StatusOK, response)
}

func (c *userController) CheckIn(context *gin.Context) {

	// Get current user from auth middleware
	currentUser := context.MustGet("currentUser").(entity.User)

	// Checkin
	result, err := c.userService.CreateAttendance(currentUser.Id, "check in")
	if err != nil {
		// Build response error
		response := helper.BuildErrorResponse("failed to process request", "check in error", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//Build response if success
	response := helper.BuildResponse(true, "successfully check in!", helper.CreateAttendanceResponse(result))
	context.JSON(http.StatusOK, response)
}

func (c *userController) CheckOut(context *gin.Context) {
	// Get current user from auth middleware
	currentUser := context.MustGet("currentUser").(entity.User)

	// Checkout
	result, err := c.userService.CreateAttendance(currentUser.Id, "check out")
	if err != nil {
		// Build response error
		response := helper.BuildErrorResponse("failed to process request", "check out error", helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//Build response if success
	response := helper.BuildResponse(true, "successfully check out!", helper.CreateAttendanceResponse(result))
	context.JSON(http.StatusOK, response)
}

func (c *userController) CreateActivity(context *gin.Context) {
	// Get current user from auth middleware
	currentUser := context.MustGet("currentUser").(entity.User)

	// Fill the activity for validation
	var activity entity.Activity
	errDTO := context.ShouldBind(&activity)
	if errDTO != nil {
		response := helper.BuildErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Create activity
	result, err := c.userService.CreateActivity(activity, currentUser.Id)
	if err != nil {
		// Build response error
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//Build response if success
	response := helper.BuildResponse(true, "successfully created activity!", helper.CreateActivityResponse(result))
	context.JSON(http.StatusCreated, response)
}

func (c *userController) UpdateActivity(context *gin.Context) {
	// Get activity data
	act_id := context.Param("id_activity")

	// Fill the updatedData for validation
	var updatedData entity.Activity
	errDTO := context.ShouldBind(&updatedData)
	if errDTO != nil {
		response := helper.BuildErrorResponse("failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusUnprocessableEntity, response)
		return
	}

	// Update activity
	result, err := c.userService.UpdateActivity(updatedData, act_id)
	if err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	//Build response if success
	response := helper.BuildResponse(true, "successfully update activity!", helper.CreateActivityResponse(result))
	context.JSON(http.StatusCreated, response)
}

func (c *userController) DeleteActivity(context *gin.Context) {
	// Get activity data
	act_id := context.Param("id_activity")

	// Delete
	err := c.userService.DeleteActivity(act_id)
	if err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Build response if success
	res := helper.BuildResponse(true, "activity deleted!", helper.EmptyObj{})
	context.JSON(http.StatusOK, res)
}

func (c *userController) GetAttendancesHistory(context *gin.Context) {
	// Get current user from auth middleware
	currentUser := context.MustGet("currentUser").(entity.User)

	// Get attendance history
	attendances, err := c.userService.GetAllAttendances(currentUser.Id)
	if err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Build response if success
	res := helper.BuildResponse(true, "successfully get attendances history!", helper.CreateAttendanceResponses(attendances))
	context.JSON(http.StatusOK, res)
}

func (c *userController) GetActivityHistoryByDate(context *gin.Context) {

	// Take start date and end date from querry
	startDate := helper.StringToUnixMilli(context.Query("startDate"))
	endDate := helper.StringToUnixMilli(context.Query("endDate"))

	// Get activity history
	activities, err := c.userService.GetActivityHistoryByDate(startDate, endDate)
	if err != nil {
		response := helper.BuildErrorResponse("failed to process request", err.Error(), helper.EmptyObj{})
		context.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Build response if success
	res := helper.BuildResponse(true, "successfully get activity history!", helper.CreateActivityResponses(activities))
	context.JSON(http.StatusOK, res)
}
