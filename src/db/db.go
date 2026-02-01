package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Conn *sql.DB

func Init(path string){
	var err error
	Conn, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal("No se pudo abrir la db: ", err)
	}

  createTables := `
    CREATE TABLE IF NOT EXISTS proyectos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        nombre TEXT NOT NULL UNIQUE,
        url_general_html TEXT NOT NULL,
        url_detallada_html TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS archivos_src (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        proyecto_id INTEGER NOT NULL,
        nombre TEXT NOT NULL,
        download_url TEXT NOT NULL,
        FOREIGN KEY(proyecto_id) REFERENCES proyectos(id) ON DELETE CASCADE
    );
    `

		_, err = Conn.Exec(createTables)
		if err != nil {
			log.Fatal("No se pudieron crear las tablas: ", err)
		}
}

func leerProyecto
