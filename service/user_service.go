package service

import (
	"armiariyan/attendances-system-jwt/dto"
	"armiariyan/attendances-system-jwt/entity"
	"armiariyan/attendances-system-jwt/helper"
	"armiariyan/attendances-system-jwt/repository"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(userInput dto.RegisterDTO) (entity.User, error)
	VerifyCredential(userInput dto.LoginDTO) (entity.User, error)
	GetUserById(user_id int) (entity.User, error)

	CreateAttendance(user_id int, label string) (entity.Attendance, error)

	CreateActivity(input entity.Activity, user_id int) (entity.Activity, error)
	UpdateActivity(input entity.Activity, act_id string) (entity.Activity, error)
	DeleteActivity(act_id string) error

	GetAllAttendances(user_id int) ([]entity.Attendance, error)
	GetActivityHistoryByDate(startDate, endDate int64) ([]entity.Activity, error)

	IsCheckInToday(user_id int) (bool, error)
	IsDuplicateEmail(email string) (bool, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		userRepository: repository,
	}
}

func (service *userService) CreateUser(userInput dto.RegisterDTO) (entity.User, error) {
	// Mapping into entity user for input database using library smapping
	fmt.Println("userInput -->", userInput)
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&userInput))
	if err != nil {
		log.Fatalf("Failed map %v", err)
	}

	// Hash Password
	userToCreate.Password = helper.HashAndSalt([]byte(userToCreate.Password))

	fmt.Println("userToCreate -->", userToCreate)
	// Check Duplicate Email
	ok, err := service.IsDuplicateEmail(userToCreate.Email)
	fmt.Println("ok -->", ok)
	if err != nil {
		return userToCreate, err
	}

	if ok {
		return userToCreate, errors.New("email has been used")
	}

	// Save
	user, err := service.userRepository.Save(userToCreate)
	if err != nil {
		return user, err
	}

	// Success
	return user, nil
}

func (service *userService) IsDuplicateEmail(email string) (bool, error) {
	user, err := service.userRepository.GetUserByEmail(email)
	// If user with given email founded
	if user.Id != 0 {
		return true, err
	}

	return false, nil
}

func (service *userService) VerifyCredential(userInput dto.LoginDTO) (entity.User, error) {
	// Get Data By Email
	user, err := service.userRepository.GetUserByEmail(userInput.Email)
	if err != nil {
		return user, err
	}

	fmt.Println("userInput.Password -->", userInput.Password)
	fmt.Println("user.Password -->", user.Password)

	// Compare Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (service *userService) GetUserById(user_id int) (entity.User, error) {
	return service.userRepository.GetUserById(user_id)
}

func (service *userService) CreateAttendance(user_id int, label string) (entity.Attendance, error) {
	// Make Attendance Data
	attendance := entity.Attendance{
		Id:     helper.GenerateIdAttendance(),
		UserId: user_id,
		Label:  label,
		Date:   time.Now().UnixMilli(),
		Time:   time.Now().UnixMilli(),
	}

	return service.userRepository.SaveAttendance(attendance)
}

func (service *userService) CreateActivity(input entity.Activity, user_id int) (entity.Activity, error) {
	// Create activity data
	activity := entity.Activity{
		Id:          helper.GenerateIdActivity(),
		UserId:      user_id,
		Description: input.Description,
		DateCreated: time.Now().UnixMilli(),
		TimeCreated: time.Now().UnixMilli(),
	}
	return service.userRepository.SaveActivity(activity)
}

func (service *userService) UpdateActivity(input entity.Activity, act_id string) (entity.Activity, error) {
	// Check if activity exist
	activity, err := service.userRepository.GetActivityById(act_id)

	if err != nil {
		return activity, err
	}

	if helper.IsActivityEmpty(activity) {
		//Build response error because activity data empty
		return activity, errors.New("activity with given id not found")
	}

	// Create new activity data
	activity = entity.Activity{
		Id:          activity.Id,
		UserId:      activity.UserId,
		Description: input.Description,
		DateCreated: activity.DateCreated,
		TimeCreated: activity.TimeCreated,
	}

	// Update
	activity, err = service.userRepository.UpdateActivity(activity)
	if err != nil {
		return activity, err
	}

	return activity, nil
}

func (service *userService) DeleteActivity(act_id string) error {
	// Check if activity exist
	activity, err := service.userRepository.GetActivityById(act_id)
	if err != nil {
		return err
	}

	if helper.IsActivityEmpty(activity) {
		//Build response error because activity data empty
		return errors.New("activity with given id not found")

	}

	err = service.userRepository.DeleteActivity(activity)
	if err != nil {
		return err
	}

	return nil
}

func (service *userService) GetActivityHistoryByDate(startDate, endDate int64) ([]entity.Activity, error) {
	activities, err := service.userRepository.GetActivityHistoryByDate(startDate, endDate)
	if err != nil {
		return activities, err
	}

	// Check if activity in range date input empty
	if helper.IsActivitiesEmpty(activities) {
		return activities, errors.New("activities in that range date is empty")
	}

	return activities, nil
}

func (service *userService) GetAllAttendances(user_id int) ([]entity.Attendance, error) {
	attendances, err := service.userRepository.GetAttendancesHistory(user_id)
	if err != nil {
		return attendances, err
	}

	// If attendances history empty
	if len(attendances) == 0 {
		return attendances, errors.New("attendances history is empty")
	}

	return attendances, nil
}

func (service *userService) IsCheckInToday(user_id int) (bool, error) {
	attendances, err := service.userRepository.GetAttendancesHistory(user_id)
	if err != nil {
		return false, err
	}

	if !helper.IsCheckIn(attendances) {
		return false, errors.New("please check in first")
	}

	return true, nil
}
