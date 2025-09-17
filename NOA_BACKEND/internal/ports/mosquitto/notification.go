package mosquitto

import (
	"encoding/json"
	"log"
)

type NotificationPayload struct {
	Message string `json:"message"`
}

// HandleNotificationMQTT publishes a notification to the specified MQTT topic
func HandleNotificationMQTT(topic string) {
	if client == nil {
		log.Println("MQTT client is not initialized")
		return
	}

	log.Printf("Publish to MQTT topic: /notification/%s\n", topic)

	// Create a JSON payload
	payload := NotificationPayload{
		Message: "Notification from Gyro Server",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error creating JSON payload:", err)
		return
	}

	// Publish the notification to the specified topic
	if token := client.Publish("/notification/"+topic, 0, false, payloadBytes); token.Wait() && token.Error() != nil {
		log.Println("Error publishing notification:", token.Error())
	} else {
		log.Println("Notification published successfully")
	}
}
