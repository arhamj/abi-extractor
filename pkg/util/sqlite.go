package util

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func NewSQLiteDB(dbFilePath, migrations string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(migrations)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("migrations successful!")
	return db, nil
}
