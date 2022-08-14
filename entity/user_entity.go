package entity

type User struct {
	Id         int          `gorm:"primary_key:auto_increment" json:"id"`
	Name       string       `gorm:"type:varchar(128)" json:"name"`
	Email      string       `gorm:"type:varchar(128)" json:"email"`
	Password   string       `gorm:"type:varchar(255)" json:"-"`
	Activity   []Activity   `json:"-"`
	Attendance []Attendance `json:"-"`
	Token      string       `gorm:"-" json:"token,omitempty"`
}
