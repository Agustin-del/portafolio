package main

import (
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"net/smtp"
	"os"
	"portafolio/ui/pages"
)

const csp = "default-src 'none'; script-src https://cdn.jsdelivr.net; style-src 'self'; img-src 'self'; connect-src 'self';"

type Proyecto struct {
	Titulo                string
	Descripcion_general   string
	Descripcion_detallada string
	fuentes               []string
}

func newMail(de, asunto, mensaje string) pages.Mail {
	return pages.Mail{
		De:      de,
		Asunto:  asunto,
		Mensaje: mensaje,
	}
}

func enviarEmail(c pages.Mail) error {
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
		c.De,
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

func render(c echo.Context, component templ.Component) {
	templ.Handler(component).ServeHTTP(
		c.Response(),
		c.Request(),
	)
}

func isHtmx(c echo.Context) bool {
	return c.Request().Header.Get("Hx-Request") == "true"
}

func main() {
	//var proyecto []Proyecto

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Hx-Request") != "true" {
				c.Response().Header().Set("Content-Security-Policy", csp)
			}

			return next(c)
		}
	})

	e.Static("/imagenes", "static/imagenes")
	e.Static("/estilos", "static/estilos")

	//TODO: quizas agregar funcionalidades para mi, metricas etc
	//	admin := e.Group("")
	e.GET("/", func(c echo.Context) error {
		if isHtmx(c) {
			return render(c, pages.InicioContenido())
		}

		return render(c, pages.Inicio())
	})

	e.GET("/contacto", func(c echo.Context) error {
		if isHtmx(c) {
			return render(c, pages.ContactoContenido(pages.Mail{}))
		}
		return render(c, pages.Contacto())
	})

	e.POST("/contacto/mail", func(c echo.Context) error {
		de := c.FormValue("email")
		asunto := c.FormValue("asunto")
		mensaje := c.FormValue("mensaje")

		//TODO: validacion minima de que esten llenos los campos
		// mantener el boton desactivado si estan vacios
		email := newMail(de, asunto, mensaje)

		if err := enviarEmail(email); err != nil {
			return render(c, pages.Contacto())
		}

		return render(c, pages.ContactoContenido(email))
	})

	if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
	}
}
