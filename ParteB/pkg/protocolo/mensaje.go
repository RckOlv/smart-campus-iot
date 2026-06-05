package protocolo


type Heartbeat struct {
	NodoID    string `json:"nodo_id"`
	Timestamp int64  `json:"timestamp"`
	Contador  int    `json:"contador"`
}

type Lectura struct {
	SensorID    string  `json:"sensor_id"`
	Temperatura float64 `json:"temperatura"`
	Timestamp   int64   `json:"timestamp"`
}

type RespuestaLectura struct {
	ID      int    `json:"id"`
	Mensaje string `json:"mensaje"`
}

type ConsultaUltimaLectura struct {
	SensorID string `json:"sensor_id"`
}