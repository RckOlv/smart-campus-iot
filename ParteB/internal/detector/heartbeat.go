package detector

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"sd-comunicacion/pkg/protocolo"
)

// ==========================================
// 1. EL ENVIADOR (Corre en el Servidor)
// ==========================================

type Enviador struct {
	destino   string
	intervalo time.Duration
	nodoID    string
	contador  int
}

func NuevoEnviador(destino string, intervalo time.Duration, nodoID string) *Enviador {
	return &Enviador{
		destino:   destino,
		intervalo: intervalo,
		nodoID:    nodoID,
		contador:  0,
	}
}

func (e *Enviador) Iniciar() {
	destinos := strings.Split(e.destino, ",")
	ticker := time.NewTicker(e.intervalo)
	defer ticker.Stop()

	fmt.Printf("[HEARTBEAT-TX] Motor de latidos iniciado para destinos: %s\n", e.destino)

	for range ticker.C {
		e.contador++
		latido := protocolo.Heartbeat{
			NodoID:    e.nodoID,
			Timestamp: time.Now().Unix(),
			Contador:  e.contador,
		}
		data, _ := json.Marshal(latido)

		for _, d := range destinos {
			target := strings.TrimSpace(d)
			addr, err := net.ResolveUDPAddr("udp", target)
			if err != nil {
				continue
			}

			conn, err := net.DialUDP("udp", nil, addr)
			if err != nil {
				continue
			}
			conn.Write(data)
			conn.Close() 
		}
	}
}

// ==========================================
// 2. EL RECEPTOR (Corre en el Cliente)
// ==========================================

type Receptor struct {
	puerto  string
	timeout time.Duration
	ultimo  time.Time
	estado  string
}

func NuevoReceptor(puerto string, timeout time.Duration) *Receptor {
	return &Receptor{
		puerto:  puerto,
		timeout: timeout,
		ultimo:  time.Now(),
		estado:  "unknown",
	}
}

func (r *Receptor) Escuchar() {
	addr, err := net.ResolveUDPAddr("udp", r.puerto)
	if err != nil {
		fmt.Printf("[HEARTBEAT-RX] Error resolviendo puerto: %v\n", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("[HEARTBEAT-RX] Error escuchando en UDP: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("[HEARTBEAT-RX] Escuchando latidos en el puerto %s\n", r.puerto)

	go func() {
		reloj := time.NewTicker(1 * time.Second)
		defer reloj.Stop()

		for range reloj.C {
			if r.estado == "unknown" {
				continue
			}

			tiempoPasado := time.Since(r.ultimo)
			nuevoEstado := r.estado

			if tiempoPasado > 2*r.timeout {
				nuevoEstado = "dead"
			} else if tiempoPasado > r.timeout {
				nuevoEstado = "suspect"
			} else {
				nuevoEstado = "alive"
			}

			if nuevoEstado != r.estado {
				r.estado = nuevoEstado
				switch r.estado {
				case "dead":
					fmt.Printf("[ALERTA] Servidor DEAD ☠️ (hace %v que no responde)\n", tiempoPasado.Round(time.Second))
				case "suspect":
					fmt.Printf("[ALERTA] Servidor SUSPECT ⚠️ (timeout superado)\n")
				case "alive":
					fmt.Printf("[INFO] Servidor recuperado a ALIVE 🟢\n")
				}
			}
		}
	}()

	buffer := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}

		var latido protocolo.Heartbeat
		if err := json.Unmarshal(buffer[:n], &latido); err == nil {
			r.ultimo = time.Now()
			
			if r.estado == "unknown" {
				r.estado = "alive"
				fmt.Printf("[INFO] Primer latido recibido de %s. Estado: ALIVE 🟢\n", latido.NodoID)
			}
		}
	}
}