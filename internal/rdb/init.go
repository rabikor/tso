package rdb

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"treatment-scheme-organizer/internal/configs"
)

var DB *gorm.DB

func init() {
	var err error

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4",
		configs.Env.DB.User,
		configs.Env.DB.Password,
		configs.Env.DB.Host,
		configs.Env.DB.Port,
		configs.Env.DB.Name,
	)

	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Panic("Failed to connect to the Database", err)
	}
}
