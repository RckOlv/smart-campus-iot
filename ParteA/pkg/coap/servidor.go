package coap

import (
	"bytes"
	"encoding/json"
	"log"
	"sync"
	"time"

	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"

	"github.com/plgd-dev/go-coap/v3"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
)

type ServidorCoAP struct {
	sim    *sensor.Simulador
	config nodo.Configuracion
	mu     sync.RWMutex
	modo   string
}

func NuevoServidor(sim *sensor.Simulador, config nodo.Configuracion) *ServidorCoAP {
	return &ServidorCoAP{
		sim:    sim,
		config: config,
		modo:   "automatico",
	}
}

// TODO 6: Iniciar arranca el servidor UDP en el puerto 5683.
func (s *ServidorCoAP) Iniciar() {
	// 6a. Crear router (es como un policía de tránsito que dirige las peticiones)
	r := mux.NewRouter()

	// 6b. Registrar handler GET /temperatura
	r.Handle("/temperatura", mux.HandlerFunc(func(w mux.ResponseWriter, req *mux.Message) {
		if req.Code() != codes.GET {
			w.SetResponse(codes.MethodNotAllowed, message.TextPlain, bytes.NewReader([]byte("Solo se permite GET")))
			return
		}

		payload := map[string]interface{}{
			"nodo_id":     s.config.ID,
			"temperatura": s.sim.Leer(),
			"unidad":      "C",
			"timestamp":   time.Now().Format(time.RFC3339),
		}

		jsonBytes, _ := json.Marshal(payload)
		w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(jsonBytes))
		log.Printf("[CoAP] Petición GET /temperatura respondida.")
	}))

	// 6c y 6d. Registrar handlers para /config (GET y PUT)
	r.Handle("/config", mux.HandlerFunc(func(w mux.ResponseWriter, req *mux.Message) {
		
		if req.Code() == codes.GET {
			s.mu.RLock() 
			payload := map[string]interface{}{
				"modo": s.modo,
			}
			s.mu.RUnlock()

			jsonBytes, _ := json.Marshal(payload)
			w.SetResponse(codes.Content, message.AppJSON, bytes.NewReader(jsonBytes))
			log.Printf("[CoAP] Petición GET /config respondida.")
			return
		}

		if req.Code() == codes.PUT {
			body, err := req.ReadBody()
			if err != nil {
				w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("Cuerpo de petición vacío o con error")))
				return
			}

			var nuevaConfig map[string]interface{}
			if err := json.Unmarshal(body, &nuevaConfig); err != nil {
				w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte("JSON inválido")))
				return
			}

			s.mu.Lock() 
			if nuevoModo, existe := nuevaConfig["modo"].(string); existe {
				s.modo = nuevoModo
				log.Printf("[CoAP] Configuración actualizada: modo = %s", s.modo)
			}
			s.mu.Unlock()

			w.SetResponse(codes.Changed, message.TextPlain, bytes.NewReader([]byte("Configuración actualizada correctamente")))
			return
		}
		w.SetResponse(codes.MethodNotAllowed, message.TextPlain, bytes.NewReader([]byte("Método no soportado")))
	}))

	// 6e. Llamar a ListenAndServe
	log.Println("[CoAP] Servidor iniciado y escuchando en el puerto UDP :5683")
	err := coap.ListenAndServe("udp", ":5683", r)
	if err != nil {
		log.Fatalf("[CoAP Error] No se pudo iniciar el servidor: %v", err)
	}
}