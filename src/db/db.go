package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)
var Conn *sql.DB

func Init(path string)  {
	
}
