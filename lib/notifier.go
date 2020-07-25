package lib

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// NotifierChan - channel to send messages to MQ
var NotifierChan = make(chan string, 1000)

func mqttConnect(hostIP string) (mqtt.Client, error) {
	opts := createClientOptions(hostIP)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		//log.Println(token.Error())
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

// NotifierService function sends messages to MQ broker
func NotifierService(notifierChan chan string, hostIP string) {
	log.Println("Trying to Connect to MQ Broker....")
	clientObj, connectErr := mqttConnect(hostIP)
	if connectErr != nil {
		log.Println("Could not connect MQ Broker: ", connectErr)
		return
	}

	log.Println("Connected to Image Store Notification Service")
	for {
		tempMSG, ok := <-notifierChan
		if !ok {
			log.Println("Notifier channel closed, notification service interupted")
			mqttCloseConnection(clientObj)
			return
		}
		token := clientObj.Publish("imagestore/notifier", 1, false, tempMSG)
		token.Wait()
	}
}
