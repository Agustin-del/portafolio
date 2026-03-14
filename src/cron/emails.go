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

type Emisor struct {
	from string
	pass string
}

func NewEmisor() *Emisor {
	return &Emisor{
		from: os.Getenv("USUARIO_SMTP"),
		pass: os.Getenv("PASS_SMTP"),
	}
}

func (e *Emisor) Enviar(m db.Mensaje) error {
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

func ProcesarPendientes(emisor *Emisor) {
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
		if err := emisor.Enviar(m); err != nil {
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
	emisor := NewEmisor()

	ticker := time.NewTicker(IntervaloMinutos * time.Minute)

	log.Printf("Cron de emails iniciado (cada %d minutos)", IntervaloMinutos)

	ProcesarPendientes(emisor)

	for range ticker.C {
		ProcesarPendientes(emisor)
	}
}
