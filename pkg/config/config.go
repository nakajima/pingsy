package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Configuration struct {
	URLs []string          `json:"urls"`
	MQTT MQTTConfiguration `json:"mqtt"`
}

type MQTTConfiguration struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (configuration *MQTTConfiguration) GetHost() string {
	return configuration.Host
}

func (configuration *MQTTConfiguration) GetPort() int {
	if configuration.Port == 0 {
		return 1883
	}

	return configuration.Port
}

func (configuration *MQTTConfiguration) GetUsername() string {
	return configuration.Username
}

func (configuration *MQTTConfiguration) GetPassword() string {
	return configuration.Password
}

var config Configuration

func init() {
	// Open our jsonFile
	jsonFile, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var configuration Configuration

	json.Unmarshal(byteValue, &configuration)
	fmt.Printf("Config is %+v\n", configuration)
	config = configuration
}

func URLS() []string {
	return config.URLs
}

func MQTT() *MQTTConfiguration {
	return &config.MQTT
}
