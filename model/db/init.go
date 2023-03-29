package db

import (
	"os"

	"github.com/Kidsunbo/kie_toolbox_go/logs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func MustInit() {
	if err := InitUserDB(); err != nil {
		logs.Fatal("failed to init UserDB, %v", err)
	}
}

func InitUserDB() error {
	dsn := os.Getenv("DB_POSTGRES_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}
