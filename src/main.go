package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"portafolio/cron"
	"portafolio/db"
	"portafolio/ui/pages"
)

const csp = "default-src 'none'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self'; img-src 'self'; connect-src 'self';"

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
	if err := db.Init("data/portafolio.db"); err != nil {
		log.Fatalf("Error inicializando base de datos: %v", err)
	}
	defer db.Close()

	go cron.IniciarCron()

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
			if !isHtmx(c) {
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
		render(c, pages.Proyectos(), http.StatusOK)
		return nil
	})

	e.GET("/contacto", func(c echo.Context) error {
		render(c, pages.Contacto(), http.StatusOK)
		return nil
	})

	e.POST("/contacto/mail", func(c echo.Context) error {
		ip := c.RealIP()

		if c.FormValue("website_url") != "" {
			c.Logger().Warn("Honey pot activado desde IP: " + ip)
			render(c, pages.ContactoExito(), http.StatusOK)
			return nil
		}

		hora := 60
		if count, err := db.ContarMensajesPorIP(ip, hora); err == nil && count >= 3 {
			c.Logger().Warn("Límite de mensajes excedido para IP: " + ip)
			render(c, pages.ContactoExito(), http.StatusOK)
			return nil
		}

		de := c.FormValue("email")
		asunto := c.FormValue("asunto")
		mensaje := c.FormValue("mensaje")

		if asunto == "" || mensaje == "" {
			render(c, pages.ContactoVacio(), http.StatusUnprocessableEntity)
			return nil
		}

		if _, err := db.GuardarMensaje(de, asunto, mensaje, ip); err != nil {
			c.Logger().Error("Error guardando mensaje", "error", err)
			render(c, pages.ContactoError(), http.StatusInternalServerError)
			return nil
		}

		render(c, pages.ContactoExito(), http.StatusOK)
		return nil
	})

	go func() {
		if err := e.Start(":42069"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("Fallo al intentar iniciar el servidor", "error", err)
		}
	}()

	signalCh := make(chan os.Signal, 1)	
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	<- signalCh 
	log.Printf("Señal, cerrando servidor...")
	e.Close()
	log.Println("Servidor cerrado")
}
