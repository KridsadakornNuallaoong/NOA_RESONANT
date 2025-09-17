package mosquitto

import (
	"encoding/json"
	"log"

	"GOLANG_SERVER/configs/env"
	"GOLANG_SERVER/internal/adapters/db"
	schema "GOLANG_SERVER/pkg/schema"

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
<<<<<<< Updated upstream:NOA_BACKEND/components/protocal/mosquitto/mq.go
		// log.Printf("Sub topic: %s\n", msg.Topic())

=======
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
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/mosquitto/mq.go
		// Process the message and store it in the database
		var data schema.GyroData
		// Unmarshal the JSON message into the GyroData struct
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Println("Error unmarshaling message:", err)
			return
		}
<<<<<<< Updated upstream:NOA_BACKEND/components/protocal/mosquitto/mq.go
=======
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

		log.Println("UserId: ", user.ID)
		HandleNotificationMQTT(string(user.ID)) // Send notification to MQTT topic

>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/mosquitto/mq.go
		if data != (schema.GyroData{}) { // if received data successfully, log it
			// log.Println("Received data from MQTT from topic:", msg.Topic())
			// log.Println("From device:", data.DeviceAddress)

			// Log the received data
			// log.Printf("Received data from MQTT topic '%s'\n", msg.Topic())

<<<<<<< Updated upstream:NOA_BACKEND/components/protocal/mosquitto/mq.go
=======
			// log.Println("Received data from MQTT from topic:", msg.Topic())
			// log.Println("From device:", data.DeviceAddress)

			// Log the received data
			// log.Printf("Received data from MQTT topic '%s'\n", msg.Topic())

>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/mosquitto/mq.go
			// Store the data in the database
			if _, err := db.StoreGyroData(data); err != nil {
				log.Println("Error storing data in database:", err)
			}
<<<<<<< Updated upstream:NOA_BACKEND/components/protocal/mosquitto/mq.go
			PublishMQTT("MQTT_TOPIC", string(msg.Payload())) // Publish the device address to the topic
=======
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/mosquitto/mq.go
		}

	}); token.Wait() && token.Error() != nil {
		log.Fatal("Error subscribing to topic:", token.Error())
	}

<<<<<<< Updated upstream:NOA_BACKEND/components/protocal/mosquitto/mq.go
	log.Println("MQTT client ready to connect and subscribe to topic.")
}

// PublishMQTT publishes a message to the specified topic
func PublishMQTT(topic string, message string) error {
	// Publish the message to the topic
	if token := client.Publish(topic, 0, false, message); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
=======
	// Log the successful connection and subscription
	log.Println("MQTT client ready to connect and subscribe to topic.")
}
>>>>>>> Stashed changes:NOA_BACKEND/internal/ports/mosquitto/mq.go
