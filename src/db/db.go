package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Conn *sql.DB

func Init(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	Conn = db

	schema := `
	CREATE TABLE IF NOT EXISTS mensajes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		de TEXT,
		asunto TEXT,
		mensaje TEXT,
		ip TEXT,
		estado TEXT DEFAULT 'pendiente',
		intentos INTEGER DEFAULT 0,
		ultimo_error TEXT,
		creado_en DATETIME DEFAULT CURRENT_TIMESTAMP,
		enviado_en DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_estado ON mensajes(estado);
	`

	if _, err := Conn.Exec(schema); err != nil {
		return err
	}

	log.Println("Base de datos inicializada correctamente")
	return nil
}

func Close() error {
	if Conn != nil {
		return Conn.Close()
	}
	return nil
}
