package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/plgd-dev/go-coap/v3/udp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := udp.Dial("127.0.0.1:15683")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	resp, err := conn.Get(ctx, "/temperatura")
	if err != nil {
		log.Fatal(err)
	}

	body, err := resp.ReadBody()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Respuesta CoAP:")
	fmt.Println(string(body))
}