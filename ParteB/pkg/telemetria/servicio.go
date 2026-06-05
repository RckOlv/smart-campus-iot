package telemetria

import (
	"fmt"
	"sd-comunicacion/pkg/protocolo"
	"sync"
)


type Telemetria struct {
	mu          sync.RWMutex
	lecturas    map[string]protocolo.Lectura
	contadorIDs int
}

func NuevaTelemetria() *Telemetria {
	return &Telemetria{
		lecturas:    make(map[string]protocolo.Lectura),
		contadorIDs: 0,
	}
}

func (t *Telemetria) RegistrarLectura(args protocolo.Lectura, resp *protocolo.RespuestaLectura) error {
	// Bloqueamos la base de datos porque vamos a escribir
	t.mu.Lock()
	defer t.mu.Unlock()

	t.lecturas[args.SensorID] = args
	
	t.contadorIDs++

	resp.ID = t.contadorIDs
	resp.Mensaje = fmt.Sprintf("Lectura de %s (%.2f°C) registrada exitosamente", args.SensorID, args.Temperatura)

	fmt.Printf("[SERVIDOR-RPC] 📥 Recibido: %s reportó %.2f°C\n", args.SensorID, args.Temperatura)
	

	return nil
}

func (t *Telemetria) ObtenerUltimaLectura(args protocolo.ConsultaUltimaLectura, resp *protocolo.Lectura) error {

	t.mu.RLock()
	defer t.mu.RUnlock()

	lectura, existe := t.lecturas[args.SensorID]
	if !existe {
		return fmt.Errorf("no hay registros para el sensor: %s", args.SensorID)
	}

	*resp = lectura
	
	fmt.Printf("[SERVIDOR-RPC] 📤 Consulta respondida para: %s\n", args.SensorID)
	return nil
}