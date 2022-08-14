package entity

type Activity struct {
	Id          string `gorm:"primaryKey;type:varchar(128)" json:"id"`
	UserId      int    `json:"id_user"`
	Description string `gorm:"type:varchar(128)" json:"description"`
	DateCreated int64  `json:"date_created"`
	TimeCreated int64  `json:"time_created"`
	User        User   `gorm:"foreignKey:UserId" json:"-"`
}
