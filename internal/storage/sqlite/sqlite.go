package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/saurabhdagwar/students-api/internal/config"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqllite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	name TEXT,
	email TEXT,
	age INTERGER
	)`)
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		Db: db,
	}, nil
}
