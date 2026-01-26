package main

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"portafolio/ui/pages"
)

type Proyecto struct{
	Titulo string
	Descripcion_general string
	Descripcion_detallada string
	fuentes []string
}

type Contacto struct {
	Email string
	Subject string
	Message string
}

func noCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-store", "no-cache", "must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")
		return next(c)
	}
}
func isHtmx(c echo.Context) bool {
	return c.Request().Header.Get("Hx-Request") == "true"
}

func render(c echo.Context, component templ.Component) error {

	templ.Handler(component).ServeHTTP(
		c.Response(),
		c.Request(),
	)

	return nil
}

func newContacto (email, subject, message string) Contacto{
	return Contacto{
		Email: email,
		Subject: subject,
		Message: message,
	}
}

func main() {
	//var proyecto []Proyecto

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Static("/imagenes", "static/imagenes")
	e.Static("/estilos", "static/estilos")

	e.GET("/", func(c echo.Context) error {
		if isHtmx(c) {
			return render(c, pages.InicioContenido())
		}

		return render(c, pages.Inicio())
	})

	e.GET("/contacto", func(c echo.Context) error {
		//TODO:validar datos	
		/*
		email := c.FormValue("email")
		subject := c.FormValue("subject")
		message := c.FormValue("message")
		contacto := newContacto(email, subject, message)

		*/
		//TODO: enviarlo a gmail
		if isHtmx(c) {
			return render(c, pages.ContactoContenido())
		}
		return render(c, pages.Contacto())
	})

	if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
	}
}
