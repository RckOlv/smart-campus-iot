package detector

// TODO 5-8: Implementar envio y recepcion de heartbeats UDP.
// Necesitaras importar:
//   "encoding/json"
//   "fmt"
//   "net"
//   "time"
//   "sd-comunicacion/pkg/protocolo"

// Enviador se encarga de enviar heartbeats UDP periodicamente
type Enviador struct {
	destino   string
	intervalo time.Duration // TODO: usar time.Duration en vez de int64
	nodoID    string
	contador  int
}

// TODO 5: Implementar la funcion NuevaEnviador.
// Debe recibir destino (string), intervalo (time.Duration) y nodoID (string).

// TODO 6: Implementar el metodo (e *Enviador) Iniciar().
// Debe enviar Heartbeat cada 'intervalo' por UDP al destino configurado.

// Receptor escucha heartbeats y detecta si dejan de llegar.
// Debe manejar estados: alive -> suspect -> dead.
type Receptor struct {
	puerto  string
	timeout time.Duration // TODO: usar time.Duration en vez de int64
	// ultimo debe guardar time.Time o timestamp del ultimo heartbeat recibido
	ultimo time.Time
	// estado puede ser "alive", "suspect" o "dead"
	activo bool
}

// TODO 7: Implementar la funcion NuevoReceptor.
// Debe recibir puerto (string) y timeout (time.Duration).

// TODO 8: Implementar el metodo (r *Receptor) Escuchar().
// Debe:
//   - Escuchar UDP en 'puerto'
//   - Decodificar mensajes JSON tipo protocolo.Heartbeat
//   - Actualizar ultimo timestamp al recibir
//   - En una goroutine separada, revisar periodicamente:
//       si time.Since(ultimo) > timeout: pasar a "suspect"
//       (opcional) si time.Since(ultimo) > 2*timeout: pasar a "dead"
//   - Imprimir cambios de estado por consola
