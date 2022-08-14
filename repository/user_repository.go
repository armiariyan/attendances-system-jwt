package repository

import (
	"armiariyan/attendances-system-jwt/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Save(user entity.User) (entity.User, error)
	VerifyCredential(email string) (entity.User, error)

	GetUserById(user_id int) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetActivityById(act_id string) (entity.Activity, error)

	SaveAttendance(attendance entity.Attendance) (entity.Attendance, error)

	SaveActivity(activity entity.Activity) (entity.Activity, error)
	UpdateActivity(activity entity.Activity) (entity.Activity, error)
	DeleteActivity(activity entity.Activity) error

	GetAttendancesHistory(user_id int) ([]entity.Attendance, error)
	GetActivityHistoryByDate(startDate, endDate int64) ([]entity.Activity, error)
}

type userConnection struct {
	connection *gorm.DB
}

// Construct
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) Save(user entity.User) (entity.User, error) {
	err := db.connection.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, err
}

func (db *userConnection) VerifyCredential(email string) (entity.User, error) {
	var user entity.User

	err := db.connection.Where("email = ?", email).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (db *userConnection) GetUserById(user_id int) (entity.User, error) {
	var user entity.User
	err := db.connection.Where("id = ?", user_id).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (db *userConnection) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	err := db.connection.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return user, err
	}
	return user, err
}

func (db *userConnection) GetActivityById(act_id string) (entity.Activity, error) {
	var activity entity.Activity
	err := db.connection.Where("id = ?", act_id).Take(&activity).Error
	if err != nil {
		return activity, err
	}
	return activity, err
}

func (db *userConnection) SaveAttendance(activity entity.Attendance) (entity.Attendance, error) {
	err := db.connection.Save(&activity).Error
	if err != nil {
		return activity, err
	}
	return activity, err
}

func (db *userConnection) SaveActivity(activity entity.Activity) (entity.Activity, error) {
	err := db.connection.Create(&activity).Error
	if err != nil {
		return activity, err
	}
	return activity, err
}

func (db *userConnection) UpdateActivity(activity entity.Activity) (entity.Activity, error) {
	err := db.connection.Save(&activity).Error

	if err != nil {
		return activity, err
	}

	return activity, nil
}

func (db *userConnection) DeleteActivity(activity entity.Activity) error {
	err := db.connection.Delete(&activity).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *userConnection) GetActivityHistoryByDate(startDate, endDate int64) ([]entity.Activity, error) {
	var activities []entity.Activity
	err := db.connection.Where("date_created >= ? AND date_created <= ?", startDate, endDate).Find(&activities).Error
	if err != nil {
		return activities, err
	}

	return activities, nil
}

func (db *userConnection) GetAttendancesHistory(user_id int) ([]entity.Attendance, error) {
	var attendances []entity.Attendance
	err := db.connection.Find(&attendances, "user_id = ?", user_id).Error
	if err != nil {
		return attendances, err
	}

	return attendances, nil
}
