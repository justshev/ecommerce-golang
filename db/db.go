package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase() (*Database, error) {
db,err := sqlx.Connect("mysql","root:root@tcp(localhost:3307)/ecomm?parseTime=true")
if err != nil {
	return nil, fmt.Errorf("Error %w", err)
}
return &Database{db: db}, nil
	
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sqlx.DB {
	return d.db
}