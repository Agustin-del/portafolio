package main

import (
	"errors"
	"html/template"
	"net/http"
	"io"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, nombre string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, nombre, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

func main() {
	e := echo.New()

	e.Static("/imagenes", "static/imagenes")
	e.Static("/estilos", "static/estilos")

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	templates := newTemplate()
	e.Renderer = templates 

	e.StdLogger.Printf("%s\n", templates.templates.DefinedTemplates())

	e.GET("/", func(c echo.Context) error {
			
		return c.Render(http.StatusOK, "layout", "sin datos")
	})

	if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
    e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
	}
}
