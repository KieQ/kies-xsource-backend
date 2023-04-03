package table

import "time"

const NameVoyage = "t_voyage"

type Voyage struct {
	ID          int32        `gorm:"column:id"`
	UserID      int32        `gorm:"column:user_id"`
	Seed        int32        `gorm:"column:seed"`
	Level       int32        `gorm:"column:level"`
	Status      VoyageStatus `gorm:"column:status"`
	Records     string       `gorm:"column:records"`
	StartTime   time.Time    `gorm:"column:start_time"`
	PassTime    time.Time    `gorm:"column:pass_time"`
	LastTryTime time.Time    `gorm:"column:last_try_time"`
}

type VoyageStatus int16

const (
	VoyageStatusInProgress VoyageStatus = iota
	VoyageStatusPass
	VoyageStatusFail
)

type VoyageRecord struct {
	CreateTime time.Time `json:"create_time"`
	Action string `json:"action"`
	Extra map[string]any `json:"extra"`
}
