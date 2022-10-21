package main

import (
	"context"
	"flag"
	"log"

	"github.com/berfarah/savant.mqtt/config"
	server "github.com/berfarah/savant.mqtt/mqtt"
	"github.com/berfarah/savant.mqtt/savant"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	flag.Parse()
	config := config.LoadConfig()

	manager, err := savant.New(config)
	if err != nil {
		log.Fatalf("Couldn't start Savant Polling Service: %v\n", err)
	}
	serv := server.New(config, manager)

	opts := mqtt.NewClientOptions()
	if config.UseSSL {
		// do tls config things
		// opts.SetTLSConfig()
	}
	if config.Username != "" {
		opts.SetUsername(config.Username)
		opts.SetPassword(config.Password)
	}

	opts.SetAutoReconnect(true)
	opts.AddBroker(config.Broker)
	opts.SetConnectionLostHandler(func(c mqtt.Client, e error) {
		log.Println("Lost connection:", e)
	})
	opts.SetClientID("savant-mqtt")
	opts.SetOnConnectHandler(serv.OnConnect)
	opts.SetDefaultPublishHandler(serv.Handler)
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if token.Wait(); token.Error() != nil {
		log.Fatalf("Couldn't connect to MQTT broker: %v\n", token.Error())
	}

	// Figure out trapping interrupt signal later
	serv.Run(context.Background())

	// Wait 200ms before exiting
	client.Disconnect(200)
}
