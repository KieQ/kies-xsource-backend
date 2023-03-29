package table

import "time"

const NameUser = "user"

type User struct {
	ID               int64     `gorm:"column:id"`
	Account          string    `gorm:"column:account"`
	Password         string    `gorm:"column:password"`
	NickName         string    `gorm:"column:nick_name"`
	Profile          string    `gorm:"column:profile"`
	Phone            string    `gorm:"column:phone"`
	Email            string    `gorm:"column:email"`
	Gender           Gender    `gorm:"column:gender"`
	SelfIntroduction string    `gorm:"column:self_introduction"`
	CreateTime       time.Time `gorm:"column:create_time"`
	UpdateTime       time.Time `gorm:"column:update_time"`
}

type Gender int32

const (
	Male Gender = iota
	Female
)
