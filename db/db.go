package db

import (
	"fmt"

	"github.com/Rakhulsr/go-form-service/config/env"
	"github.com/Rakhulsr/go-form-service/db/migrations"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDbConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset:utf8mb4&parseTime=True&loc=Local", env.ENV.DBUser, env.ENV.DBPassword, env.ENV.DBPort, env.ENV.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Errorf("Failed to open connection to DB, : %v", err)
		return nil, err
	}

	migrations.AutoMigrateModels(db)

	fmt.Println("successfully connect to the db")

	return db, nil
}
