package main

import (
	"os"
	"portafolio/db"
	"testing"
)

func TestMain(m *testing.M) {
	if err := db.Init("data/test.db"); err != nil {
		os.Exit(1)
	}
	defer db.Close()

	os.Exit(m.Run())
}

func TestGuardarYObtenerMensaje(t *testing.T) {
	id, err := db.GuardarMensaje("test@example.com", "Test", "Hola mundo", "127.0.0.1")
	if err != nil {
		t.Fatalf("Error guardando mensaje: %v", err)
	}

	if id <= 0 {
		t.Fatal("ID debería ser mayor a 0")
	}

	mensajes, err := db.ObtenerPendientes(10)
	if err != nil {
		t.Fatalf("Error obteniendo mensajes: %v", err)
	}

	if len(mensajes) == 0 {
		t.Fatal("Debería haber al menos un mensaje")
	}

	ultimo := mensajes[len(mensajes)-1]
	if ultimo.De != "test@example.com" {
		t.Errorf("Expected from: test@example.com, got: %s", ultimo.De)
	}
	if ultimo.Asunto != "Test" {
		t.Errorf("Expected subject: Test, got: %s", ultimo.Asunto)
	}
}

func TestMarcarEnviado(t *testing.T) {
	id, err := db.GuardarMensaje("test2@example.com", "Test 2", "Mensaje 2", "127.0.0.1")
	if err != nil {
		t.Fatalf("Error guardando mensaje: %v", err)
	}

	if err := db.MarcarEnviado(int(id)); err != nil {
		t.Fatalf("Error marcando enviado: %v", err)
	}

	mensajes, err := db.ObtenerPendientes(10)
	if err != nil {
		t.Fatalf("Error obteniendo mensajes: %v", err)
	}

	for _, m := range mensajes {
		if m.ID == int(id) {
			t.Error("El mensaje enviado no debería aparecer en pendientes")
		}
	}
}

func TestLimitePorIP(t *testing.T) {
	ip := "192.168.1.100"

	count, err := db.ContarMensajesPorIP(ip, 60)
	if err != nil {
		t.Fatalf("Error contando mensajes: %v", err)
	}

	if count > 0 {
		t.Errorf("No debería haber mensajes para esta IP nueva, hay: %d", count)
	}
}
