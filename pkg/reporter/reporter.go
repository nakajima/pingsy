package reporter

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nakajima/pingsy/pkg/config"
)

var client mqtt.Client

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func init() {
	configuration := config.MQTT()

	var broker = configuration.GetHost()
	var port = configuration.GetPort()
	var username = configuration.GetUsername()
	var password = configuration.GetPassword()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(username)
	opts.SetPassword(password)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client = mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

type ReportMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ReportOK(uri string) {
	parsed, err := url.Parse(uri)

	if err != nil {
		log.Fatal("Could not parse URL", uri, err.Error())
	}

	payload := ReportMessage{
		Status:  "online",
		Message: uri + " OK",
	}

	json, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Could not marshall payload: ", err.Error())
	}

	client.Publish("binary_sensor/"+parsed.Host, 0, false, json)
}

func ReportError(uri string, message string) {
	parsed, err := url.Parse(uri)

	if err != nil {
		log.Fatal("Could not parse URL", uri, err.Error())
	}

	payload := ReportMessage{
		Status:  "offline",
		Message: message,
	}

	json, err := json.Marshal(payload)

	if err != nil {
		log.Fatal("Could not marshall payload: ", err.Error())
	}

	fmt.Println("HOST IS", parsed.Hostname())

	client.Publish("uptime/"+parsed.Hostname(), 0, false, json)
}
