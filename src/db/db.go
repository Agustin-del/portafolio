package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"portafolio/ui/pages"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	githubAPIBase = "https://api.gihtub.com"
	maxRetries = 3
	retryDelay = time.Second * 2
)

var Conn *sql.DB

func Init(path string)  {
	var err error
	Conn, err = sql.Open("sqlite3", path)
	if err != nil {
		fmt.Errorf("No se pudo abrir la db: %w", err)
	}

  createTables := `
    CREATE TABLE IF NOT EXISTS proyectos (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        nombre TEXT NOT NULL UNIQUE,
        descripcion_general TEXT NOT NULL,
        descripcion_detallada TEXT NOT NULL
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
			fmt.Errorf("No se pudieron crear las tablas: ", err)
		}
}

func TraerProyectos() ([]pages.Proyecto, error) {

	baseApi := fmt.Sprintf("https://api.github.com/repos/agustin-del/proyectos/contents")
	resp, err := http.Get(baseApi)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var items []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		SubmoduleGitURL string `json:"submodule_git_url"`
	}
	
	json.Unmarshal(body, &items)

	proyectos := []pages.Proyecto{}

	for _, item := range items {
		if item.SubmoduleGitURL == "" {
			continue
		}

		repo := strings.TrimPrefix(item.SubmoduleGitURL, "https://github.com/")

		p, err := TraerProyectoRepo(repo, item.Name)
		if err == nil {
			proyectos = append(proyectos, p)
		}
	}

	return proyectos, nil
}

func TraerProyectoRepo(repo, nombre string) (pages.Proyecto, error) {
	p := pages.Proyecto{Nombre: nombre}

	api := fmt.Sprintf("https://api.github.com/repos/%s/contents", repo)

	resp, err := http.Get(api)
	if err != nil {
		return p, err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var contenidos []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		DownloadURL string `json:"download_url"`
		URL string `json:"url"`
	}

	json.Unmarshal(body, &contenidos)

	for _, c := range contenidos {
		switch c.Name {
		case "descripcion-general.html":
			p.DescripcionGeneral = leerArchivo(c.DownloadURL)
		case "descripcion-detallada.html":
			p.DescripcionDetallada = leerArchivo(c.DownloadURL)
		case "src":
			recorrerSrc(c.URL, &p)
		}
	}

	return p, nil
}

func recorrerSrc(url string, p *pages.Proyecto) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var items []struct {
		Type string `json:"type"`
		Name string `json:"name"`
		DownloadURL string `json:"download_url"`
		URL string `json:"url"`
	}

	json.Unmarshal(body, &items)

	for _, item := range items {
		if item.Type == "file" {
			contenido := leerArchivo(item.DownloadURL)
			p.Fuentes = append(p.Fuentes, contenido)
		} else if item.Type == "dir" {
			recorrerSrc(item.URL, p)
		}
	}
}

func leerArchivo(url string) string {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	return string(data)
}
