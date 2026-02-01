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

const csp = "default-src 'none'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self'; img-src 'self'; connect-src 'self';"

func newMail(de, asunto, mensaje string) pages.Mail {
	return pages.Mail{
		De:      de,
		Asunto:  asunto,
		Mensaje: mensaje,
	}
}

func enviarEmail(mail pages.Mail) error {
	de := os.Getenv("USUARIO_SMTP")
	pass := os.Getenv("PASS_SMTP")

	para := de

	auth := smtp.PlainAuth("", de, pass, "smtp.gmail.com")

	msj := fmt.Sprintf(
		"To: Contacto web <%s>\r\n"+
			"From: %s\r\n"+
			"Reply-To: %s\r\n"+
			"Subject: [Contacto] %s\r\n"+
			"\r\n"+
			"%s\r\n",
		de,
		para,
		mail.De,
		mail.Asunto,
		mail.Mensaje,
	)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		de,
		[]string{para},
		[]byte(msj),
	)
}

func render(c echo.Context, component templ.Component, code int) {
	c.Response().WriteHeader(code)
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

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		c.Response().WriteHeader(code)

		templ.Handler(
			pages.Error(code),
		).ServeHTTP(c.Response(), c.Request())
	}

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
	e.Static("/scripts", "static/scripts")

	//TODO: quizas agregar funcionalidades para mi, metricas etc
	//	admin := e.Group("")
	e.GET("/", func(c echo.Context) error {
		render(c, pages.Inicio(), http.StatusOK)

		return nil
	})

	e.GET("/proyectos", func(c echo.Context) error {
		render(c,pages.Proyectos(), http.StatusOK)
		return nil
	})

	e.GET("/contacto", func(c echo.Context) error {
		render(c, pages.Contacto(), http.StatusOK)
		return nil
	})

	e.POST("/contacto/mail", func(c echo.Context) error {
		de := c.FormValue("email")
		asunto := c.FormValue("asunto")
		mensaje := c.FormValue("mensaje")

		if asunto == "" || mensaje == "" {
			render(c, pages.ContactoVacio(), http.StatusUnprocessableEntity)
			return nil
		}

		email := pages.Mail{
			De:      de,
			Asunto:  asunto,
			Mensaje: mensaje,
		}

		// mantener el boton desactivado si estan vacios
		if err := enviarEmail(email); err != nil {
			c.Logger().Error("Error enviando mail", "error", err)
			//render(c, pages.Error(http.StatusInternalServerError), http.StatusInternalServerError)
			return nil
		}

		render(c, pages.ContactoExito(), http.StatusOK)
		return nil
	})

	if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
	}
}
