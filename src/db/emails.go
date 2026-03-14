package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Mensaje struct {
	ID          int
	De          string
	Asunto      string
	Mensaje     string
	IP          string
	Estado      string
	Intentos    int
	UltimoError sql.NullString
	CreadoEn    time.Time
	EnviadoEn   sql.NullTime
}

func GuardarMensaje(de, asunto, mensaje, ip string) (int64, error) {
	result, err := Conn.Exec(
		`INSERT INTO mensajes (de, asunto, mensaje, ip, estado, intentos, creado_en) 
		 VALUES (?, ?, ?, ?, 'pendiente', 0, datetime('now', 'localtime'))`,
		de, asunto, mensaje, ip,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func ObtenerPendientes(limite int) ([]Mensaje, error) {
	rows, err := Conn.Query(
		`SELECT id, de, asunto, mensaje, ip, estado, intentos, ultimo_error, creado_en, enviado_en 
		 FROM mensajes WHERE estado = 'pendiente' AND intentos < 5 ORDER BY id ASC LIMIT ?`,
		limite,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mensajes []Mensaje
	for rows.Next() {
		var m Mensaje
		if err := rows.Scan(&m.ID, &m.De, &m.Asunto, &m.Mensaje, &m.IP, &m.Estado, &m.Intentos, &m.UltimoError, &m.CreadoEn, &m.EnviadoEn); err != nil {
			return nil, err
		}
		mensajes = append(mensajes, m)
	}
	return mensajes, nil
}

func MarcarEnviado(id int) error {
	_, err := Conn.Exec(
		`UPDATE mensajes SET estado = 'enviado', enviado_en = datetime('now', 'localtime') WHERE id = ?`,
		id,
	)
	return err
}

func MarcarFallido(id int, errorMsg string) error {
	_, err := Conn.Exec(
		`UPDATE mensajes SET intentos = intentos + 1, ultimo_error = ?, estado = CASE WHEN intentos + 1 >= 5 THEN 'fallido' ELSE 'pendiente' END WHERE id = ?`,
		errorMsg, id,
	)
	return err
}

func ContarMensajesPorIP(ip string, ventanaMinutos int) (int, error) {
	var count int
	err := Conn.QueryRow(
		`SELECT COUNT(*) FROM mensajes WHERE ip = ? AND creado_en > datetime('now', ?)`,
		ip, "-"+fmt.Sprintf("%d minutes", ventanaMinutos),
	).Scan(&count)
	return count, err
}
