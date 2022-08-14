package entity

type Attendance struct {
	Id     string `gorm:"primaryKey;type:varchar(128)" json:"id"`
	UserId int    `json:"id_user"`
	Label  string `gorm:"type:varchar(128)" json:"label"`
	Date   int64  `json:"date"`
	Time   int64  `json:"time"`
	User   User   `gorm:"foreignKey:UserId" json:"-"`
}
