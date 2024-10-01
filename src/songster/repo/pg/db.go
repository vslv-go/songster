package pg

import (
	"github.com/pressly/goose"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

func New(pgDSL string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(pgDSL), &gorm.Config{
		TranslateError: true,
	})
	return &DB{db}, err
}

func (db *DB) RunMigrations(path string) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return goose.Up(sqlDB, path)
}
