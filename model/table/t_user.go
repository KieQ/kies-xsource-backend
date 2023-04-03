package table

import "time"

const NameUser = "t_user"

type User struct {
	ID               int32      `gorm:"column:id"`
	Phone            string     `gorm:"column:phone"`
	Email            string     `gorm:"column:email"`
	Password         string     `gorm:"column:password"`
	NickName         string     `gorm:"column:nick_name"`
	Profile          string     `gorm:"column:profile"`
	Gender           UserGender `gorm:"column:gender"`
	SelfIntroduction string     `gorm:"column:self_introduction"`
	CreateTime       time.Time  `gorm:"column:create_time"`
	UpdateTime       time.Time  `gorm:"column:update_time"`
}

type UserGender int8

const (
	Male UserGender = iota
	Female
)
