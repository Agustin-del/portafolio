package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"portafolio/cron"
	"portafolio/db"
	"portafolio/ui/pages"
)

const csp = "default-src 'none'; script-src 'self' https://cdn.jsdelivr.net; style-src 'self'; img-src 'self'; connect-src 'self';"

var proyectos map[string]map[string][]byte

func cargarProyectos() error {
	proyectos = make(map[string]map[string][]byte)
	baseDir := "proyectos"

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("error leyendo directorio proyectos: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		nombreProyecto := entry.Name()
		proyectoDir := filepath.Join(baseDir, nombreProyecto)

		proyectos[nombreProyecto] = make(map[string][]byte)

		archivos := []string{"descripcion.md", "descripcion_detallada.md"}
		for _, archivo := range archivos {
			ruta := filepath.Join(proyectoDir, archivo)
			contenido, err := os.ReadFile(ruta)
			if err != nil {
				log.Printf("Advertencia: no se pudo leer %s: %v", ruta, err)
				continue
			}
			proyectos[nombreProyecto][archivo] = contenido
		}

		srcDir := filepath.Join(proyectoDir, "src")
		filepath.WalkDir(srcDir, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				log.Printf("error leyendo entry: %v, error: %v", entry, err)
				return nil
			}

			if entry.Name() == ".git" || entry.Name() == "data" ||
				entry.Name() == "vendor" || entry.Name() == "tmp" || 
				entry.Name() == "proyectos"{
				return filepath.SkipDir
			}

			entryInfo, err := entry.Info()
			if err != nil {
				log.Printf("error obteniendo informacion del archivo: %v", err)
				return nil
			}

			if !strings.HasSuffix(entryInfo.Name(), ".go") &&
				!strings.HasSuffix(entryInfo.Name(), ".mod") &&
				!strings.HasSuffix(entryInfo.Name(), ".sum") &&
				!strings.HasSuffix(entryInfo.Name(), ".js") &&
				!strings.HasSuffix(entryInfo.Name(), ".css") &&
				!strings.HasSuffix(entryInfo.Name(), ".html") &&
				!strings.HasSuffix(entryInfo.Name(), ".templ") {
				return nil
			}

			if entryInfo.Size() > 1<<20 {
				log.Printf("tamaño gigante: %v", entryInfo.Name())
				return nil
			}

			contenido, err := os.ReadFile(path)
			if err != nil {
				log.Printf("error leyendo archivo: %v", path)
				return nil
			}

			relPath, _ := filepath.Rel(proyectoDir, path)
			proyectos[nombreProyecto][relPath] = contenido
			return nil
		})
	}

	log.Printf("Proyectos cargados: %d", len(proyectos))
	return nil
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
	if err := cargarProyectos(); err != nil {
		log.Fatalf("Error cargando proyectos: %v", err)
	}

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
    nombres := make([]string, 0, len(proyectos))
    for nombre := range proyectos {
      nombres = append(nombres, nombre)
    }
    render(c, pages.Proyectos(nombres), http.StatusOK)
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
	<-signalCh
	log.Printf("Señal, cerrando servidor...")
	//TODO:estoy cerrando el servidor, como esta funcion salga va a ejecutar el defer de la db es decir va a cerrar el archivo de sqlite
	//la cosa me parece que no estoy cerrando el ticker, es decir si el cron estaba ejecutando me parece que corta abruptamente.
	//FIX: acaso close o shutdown deberia usar?
	e.Close()
	log.Println("Servidor cerrado")
}
