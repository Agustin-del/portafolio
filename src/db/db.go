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
		id INTEGER PRIMARY KEY,
		de TEXT NOT NULL,
		asunto TEXT NOT NULL,
		mensaje TEXT NOT NULL,
		ip TEXT NOT NULL,
		estado TEXT DEFAULT 'pendiente' NOT NULL,
		intentos INTEGER DEFAULT 0 NOT NULL,
		ultimo_error TEXT NOT NULL,
		creado_en DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		enviado_en DATETIME NOT NULL
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
