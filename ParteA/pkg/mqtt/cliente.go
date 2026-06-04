package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"sd-iot/pkg/nodo"
	"sd-iot/pkg/sensor"
)

type Cliente struct {
	config   nodo.Configuracion
	interno  mqtt.Client
	opciones *mqtt.ClientOptions
}

// TODO 1: NuevoCliente crea la configuración inicial del cliente MQTT.
func NuevoCliente(config nodo.Configuracion) (*Cliente, error) {
	opts := mqtt.NewClientOptions()
	
	brokerURL := fmt.Sprintf("tcp://%s", config.BrokerMQTT)
	opts.AddBroker(brokerURL)
	opts.SetClientID(config.ID)

	// 1a y 1b: Configurar el Testamento (LWT)
	topicEstado := fmt.Sprintf("nodo/%s/estado", config.ID)
	payloadOffline := `{"estado":"offline"}`
	opts.SetWill(topicEstado, payloadOffline, 1, true)

	// 1c: Configuración extra de robustez
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	return &Cliente{
		config:   config,
		opciones: opts,
	}, nil
}

// TODO 2: Conectar establece la sesión con el broker.
func (c *Cliente) Conectar() error {
	c.interno = mqtt.NewClient(c.opciones)
	
	token := c.interno.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("falló la conexión al broker: %v", token.Error())
	}

	topicEstado := fmt.Sprintf("nodo/%s/estado", c.config.ID)
	payloadOnline := `{"estado":"online"}`
	
	tokenPub := c.interno.Publish(topicEstado, 1, true, payloadOnline)
	tokenPub.Wait()
	
	if tokenPub.Error() != nil {
		return fmt.Errorf("conectado, pero falló al publicar estado online: %v", tokenPub.Error())
	}

	log.Printf("[MQTT] Nodo %s conectado exitosamente. Estado 'online' publicado.", c.config.ID)
	return nil
}

// TODO 3: PublicarLecturas envía periódicamente las lecturas del sensor.
func (c *Cliente) PublicarLecturas(sim *sensor.Simulador, config nodo.Configuracion) {
	topic := fmt.Sprintf("campus/%s/%s/sensor/temperatura", config.Edificio, config.Aula)
	ticker := time.NewTicker(config.IntervaloSegundos)

	for t := range ticker.C {
		temperatura := sim.Leer()

		payloadMap := map[string]interface{}{
			"nodo_id":     config.ID,
			"temperatura": temperatura,
			"unidad":      "C",
			"timestamp":   t.Format(time.RFC3339),
		}

		jsonBytes, err := json.Marshal(payloadMap)
		if err != nil {
			log.Printf("[MQTT Error] No se pudo serializar JSON: %v", err)
			continue
		}

		token := c.interno.Publish(topic, 1, false, jsonBytes)
		token.Wait()

		if token.Error() != nil {
			log.Printf("[MQTT Error] Falló al publicar lectura: %v", token.Error())
		} else {
			log.Printf("[MQTT] Publicado en %s -> %s", topic, string(jsonBytes))
		}
	}
}

// TODO 4: SuscribirComandos se une al tópico de actuadores y procesa mensajes.
func (c *Cliente) SuscribirComandos(config nodo.Configuracion) error {
	topic := fmt.Sprintf("campus/%s/%s/actuador/cmd", config.Edificio, config.Aula)
	token := c.interno.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		var comando map[string]interface{}
		err := json.Unmarshal(msg.Payload(), &comando)
		if err != nil {
			log.Printf("[Actuador Error] Comando inválido recibido: %s", string(msg.Payload()))
			return
		}

		accion, _ := comando["accion"].(string)
		origen, _ := comando["origen"].(string)

		log.Printf("[Actuador] Comando recibido de '%s'. Simulando acción: >>> %s <<<", origen, accion)
	})

	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("falló la suscripción a comandos: %v", token.Error())
	}

	log.Printf("[MQTT] Nodo escuchando comandos en: %s", topic)
	return nil
}

// TODO 5: Desconectar cierra limpiamente la sesión MQTT.
func (c *Cliente) Desconectar() {
	if c.interno != nil && c.interno.IsConnected() {
		topicEstado := fmt.Sprintf("nodo/%s/estado", c.config.ID)
		payloadOffline := `{"estado":"offline"}`
		
		token := c.interno.Publish(topicEstado, 1, true, payloadOffline)
		token.Wait()

		c.interno.Disconnect(250)
		log.Printf("[MQTT] Nodo %s desconectado de forma limpia.", c.config.ID)
	}
}