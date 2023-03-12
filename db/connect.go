package db

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (gorm.DB, error) {
	dsn := os.Getenv("DSN")

	if dsn == "" {
		return gorm.DB{}, errors.New("the dsn environment variable was not found")
	}

	dsn = strings.ReplaceAll(dsn, "\n", "")

	fmt.Println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return gorm.DB{}, err
	}

	return *db, nil
}
