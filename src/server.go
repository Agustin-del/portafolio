package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"portafolio/ui/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Proyecto struct {
	Titulo                string
	Descripcion_general   string
	Descripcion_detallada string
	fuentes               []string
}

type Mail struct {
	Para   string
	Asunto  string
	Mensaje string
}

func enviarEmail(c Mail) error {
	de := os.Getenv("USUARIO_SMTP")
	pass := os.Getenv("PASS_SMTP")

	para := de

	auth := smtp.PlainAuth("", de, pass, "smtp.gmail.com")

	msj := []byte(fmt.Sprintf(
		"To: Contacto web <%s>\r\n"+
			"From: %s\r\n"+
			"Reply-To: %s\r\n"+
			"Subject: [Contacto] %s\r\n"+
			"\r\n"+
			"%s\r\n",
		de,
		para,
		c.Para,
		c.Asunto,
		c.Mensaje,
	))

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		de,
		[]string{para},
		msj,
	)
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

func newMail(para, asunto, mensaje string) Mail {
	return Mail{
		Para:   para,
		Asunto:  asunto,
		Mensaje: mensaje,
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

		if isHtmx(c) {
			return render(c, pages.ContactoContenido())
		}
		return render(c, pages.Contacto())
	})

	e.POST("/contacto", func (c echo.Context) error {
		para := c.FormValue("email")
		asunto := c.FormValue("asunto")
		mensaje := c.FormValue("mensaje")
		email := newMail(para, asunto, mensaje)

		if err := enviarEmail(email); err != nil {
			return c.String(http.StatusInternalServerError, "Error enviando email")
		}

		return c.String(http.StatusOK, "Mensaje enviado")
	})

	if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
	}
}
