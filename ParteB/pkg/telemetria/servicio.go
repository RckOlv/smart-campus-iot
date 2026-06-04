package telemetria

// TODO 1: Definir el struct Telemetria que sera el servicio RPC.
// Debe contener un mapa protegido por sync.Mutex para almacenar
// la ultima lectura de cada sensor.
// Sugerencia: usar map[string]Lectura
//
// import "sync" cuando lo necesites

type Telemetria struct {
	// COMPLETAR: mapa de lecturas y mutex
}

// TODO 2: Implementar el metodo RPC RegistrarLectura.
// Firma requerida por net/rpc:
//   func (t *Telemetria) RegistrarLectura(args Lectura, resp *RespuestaLectura) error
// Debe:
//   - Guardar la lectura en el mapa (protegiendo con mutex)
//   - Asignar un ID incremental a la respuesta
//   - Loguear la lectura recibida (import "fmt" y "time")
//   - Retornar nil en caso de exito

// TODO 3: Implementar el metodo RPC ObtenerUltimaLectura.
// Firma requerida por net/rpc:
//   func (t *Telemetria) ObtenerUltimaLectura(args ConsultaUltimaLectura, resp *Lectura) error
// Debe:
//   - Buscar en el mapa la ultima lectura del SensorID solicitado
//   - Si no existe, retornar un error con fmt.Errorf
//   - Si existe, copiar el valor a resp y retornar nil
