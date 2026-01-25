package main

import (
	"errors"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"portafolio/ui/pages"
)

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

func main() {
	e := echo.New()

	e.Static("/imagenes", "static/imagenes")
	e.Static("/estilos", "static/estilos")

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		if isHtmx(c) {
			return render(c, pages.InicioContenido())
		}

		return render(c, pages.Inicio())
	})

	e.GET("/contacto", func(c echo.Context) error {
		if isHtmx(c) {
			return render(c, pages.ContactoContenido())
		}
		return render(c, pages.Contacto())
	})

	if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
	}
}
