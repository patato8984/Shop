package db

import "database/sql"

func NewPostgesConnection(connStr string) (*sql.DB, error) {
	db, err := sql.Open("posgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
