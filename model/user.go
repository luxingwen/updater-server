package model

import "time"

type User struct {
	Id        uint      `gorm:"primaryKey" json:"id"`
	Uuid      string    `json:"uuid"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	TeamId    string    `json:"teamId"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated"`
}
