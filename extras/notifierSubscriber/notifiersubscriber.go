package main

import(
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttConnect() (mqtt.Client, error) {
	opts := createClientOptions()
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return nil, token.Error()
	}
	return client, nil
}

func createClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://mq_broker_service:1883")
	opts.SetKeepAlive(60)
	opts.SetConnectionLostHandler(func(mqttClient mqtt.Client, err error) {
		log.Printf("MQ Broker Connection lost, reason: %v\n", err)
	})
	return opts
}

func mqttCloseConnection(mqttClient mqtt.Client) {
	mqttClient.Disconnect(250)
}

func main(){
	done := make(chan bool)

	log.Println("Trying to Connect to MQ Broker....")
	clientObj, connectErr := mqttConnect()
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