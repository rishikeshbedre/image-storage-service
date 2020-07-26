package main

import (
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttConnect(hostIP string) (mqtt.Client, error) {
	opts := createClientOptions(hostIP)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return nil, token.Error()
	}
	return client, nil
}

func createClientOptions(hostIP string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://" + hostIP + ":1883")
	opts.SetKeepAlive(60)
	opts.SetConnectionLostHandler(func(mqttClient mqtt.Client, err error) {
		log.Printf("MQ Broker Connection lost, reason: %v\n", err)
	})
	return opts
}

func mqttCloseConnection(mqttClient mqtt.Client) {
	mqttClient.Disconnect(250)
}

func main() {
	hostIP := os.Getenv("HOSTIP")
	done := make(chan bool)

	log.Println("Trying to Connect to MQ Broker....")
	clientObj, connectErr := mqttConnect(hostIP)
	if connectErr != nil {
		log.Println("Could not connect MQ Broker: ", connectErr)
		return
	}

	log.Println("Subscribed to Image Store Notification Service")
	clientObj.Subscribe("imagestore/notifier", 1, func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("[%s] - %s\n", msg.Topic(), string(msg.Payload()))
	})

	<-done
}
