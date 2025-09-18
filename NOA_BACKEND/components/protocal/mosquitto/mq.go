package mosquitto

import (
	"encoding/json"
	"log"

	"GOLANG_SERVER/components/db"
	"GOLANG_SERVER/components/env"
	schema "GOLANG_SERVER/components/schema"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client

// Handle MQTT connections and messages
func HandleMQTT() {
	// Create a new MQTT client
	opts := mqtt.NewClientOptions().AddBroker(env.GetEnv("MQTT_BROKER"))
	opts.SetClientID(env.GetEnv("MQTT_CLIENT_ID"))
	opts.SetUsername(env.GetEnv("MQTT_USERNAME"))
	opts.SetPassword(env.GetEnv("MQTT_PASSWORD"))
	client = mqtt.NewClient(opts)

	// Connect to the MQTT broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("Error connecting to MQTT broker:", token.Error())
	}

	// Subscribe to the topic
	if token := client.Subscribe("vibration", 1, func(client mqtt.Client, msg mqtt.Message) {
		// Check if the message is empty
		if len(msg.Payload()) == 0 {
			log.Println("Received empty message")
			return
		}
		// Check if msg.Payload() is a valid JSON
		if !json.Valid(msg.Payload()) {
			log.Println("Received invalid JSON message:", string(msg.Payload()))
			return
		}
		// Check if msg.Payload() is a valid GyroData struct
		// Process the message and store it in the database
		var data schema.GyroData
		// Unmarshal the JSON message into the GyroData struct
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Println("Error unmarshaling message:", err)
			return
		}
		// Check userID and deviceID
		if data.UserID == "" {
			log.Println("UserID is empty")
			return
		}
		if data.DeviceID == "" {
			log.Println("DeviceID is empty")
			return
		}

		// check userID and deviceID in the database
		user, err := db.FindUserID(data.UserID)
		if err != nil {
			log.Println("Error finding user:", err)
			return
		}
		if user.ID == "" {
			log.Println("UserID not found in the database")
			return
		}
		device, err := db.FindDevice(data.DeviceID)
		if err != nil {
			log.Println("Error finding device:", err)
			return
		}
		if device.ID == "" {
			log.Println("DeviceID not found in the database")
			return
		}

		if data != (schema.GyroData{}) { // if received data successfully, log it
			// Store the data in the database
			if _, err := db.StoreGyroData(data); err != nil {
				log.Println("Error storing data in database:", err)
			}
		}

	}); token.Wait() && token.Error() != nil {
		log.Fatal("Error subscribing to topic:", token.Error())
	}

	// Log the successful connection and subscription
	log.Println("MQTT client ready to connect and subscribe to topic.")
}
