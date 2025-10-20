package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(conn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil

}
