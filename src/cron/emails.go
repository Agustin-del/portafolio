package cron

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"portafolio/db"
	"time"
)

const (
	MaxIntentos      = 5
	IntervaloMinutos = 60 
	MaxPorLote       = 10
)

type emisor struct {
	from string
	pass string
}

func newEmisor() *emisor {
	return &emisor{
		from: os.Getenv("USUARIO_SMTP"),
		pass: os.Getenv("PASS_SMTP"),
	}
}

func (e *emisor) enviar(m db.Mensaje) error {
	para := e.from

	auth := smtp.PlainAuth("", e.from, e.pass, "smtp.gmail.com")

	msg := fmt.Sprintf(
		"To: Contacto web <%s>\r\n"+
			"From: %s\r\n"+
			"Reply-To: %s\r\n"+
			"Subject: [Contacto] %s\r\n"+
			"\r\n"+
			"%s\r\n",
		para,
		para,
		m.De,
		m.Asunto,
		m.Mensaje,
	)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		e.from,
		[]string{para},
		[]byte(msg),
	)
}

func procesarPendientes(emisor *emisor) {
	mensajes, err := db.ObtenerPendientes(MaxPorLote)
	if err != nil {
		log.Printf("Error obteniendo mensajes pendientes: %v", err)
		return
	}

	if len(mensajes) == 0 {
		return
	}

	log.Printf("Procesando %d mensajes pendientes", len(mensajes))

	for _, m := range mensajes {
		if err := emisor.enviar(m); err != nil {
			log.Printf("Error enviando mensaje %d: %v", m.ID, err)
			if err := db.MarcarFallido(m.ID, err.Error()); err != nil {
				log.Printf("Error marcando mensaje %d como fallido: %v", m.ID, err)
			}
			continue
		}

		if err := db.MarcarEnviado(m.ID); err != nil {
			log.Printf("Error marcando mensaje %d como enviado: %v", m.ID, err)
			continue
		}

		log.Printf("Mensaje %d enviado exitosamente", m.ID)
	}
}

func IniciarCron() {
	emisor := newEmisor()

	ticker := time.NewTicker(IntervaloMinutos * time.Minute)

	log.Printf("Cron de emails iniciado (cada %d minutos)", IntervaloMinutos)
//acaso esta linea es necesaria? Es decir no deberia a esperar a que el ticker me mande un tick para procesar los pendientes?
//en un punto tiene sentido que procese todo cuando inicia y desde ahi cada vez que se cumple el tick, supongo
//podria chequear con una peticion a la db si hay emails nuevos o pendendientes, y sino retornar de aca, es decir cerrar el ticker
//digo para no estar corriendo siempre este cron durante todo el proceso, tendria que hacer algo en main para que cuando esa busqueda en la db
//un select count(*) from mensajes where estado=pendiente sea mayor a que se yo x corra de vuelta este cron, proceso de a 10 mensajes
//me parece poco sobre todo porque pasa concurrentemente. Acaso no seria mejor poner un minimo? Y este cron capaz dejarlo que corra,
//pero pasar la logica a la funcion que procesa para que retorne si no hay suficientes mensajes... mm no la cosa no, porque quiero que me 
//lleguen los mails. Igual el maximo me parece poco.
	procesarPendientes(emisor)

	for range ticker.C {
		procesarPendientes(emisor)
	}
}
