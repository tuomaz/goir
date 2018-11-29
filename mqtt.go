package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
)

type homeAssistantData struct {
	EventType string    `json:"event_type"`
	EventData eventData `json:"event_data"`
}

type eventData struct {
	EntityID string     `json:"entity_id"`
	OldState eventState `json:"old_state"`
	NewState eventState `json:"new_state"`
}

type eventState struct {
	EntityID    string                 `json:"entity_id"`
	State       string                 `json:"state"`
	Attributes  map[string]interface{} `json:"attributes"`
	LastChanged *time.Time             `json:"last_changed"`
	LastUpdated *time.Time             `json:"last_updated"`
}

func createAndStartMQTT(serverAddress, clientID, topic string) *client.Client {
	// Create an MQTT Client.
	cli := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("example-client"),
	})
	if err != nil {
		panic(err)
	}

	// Subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte("hass"),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: messageHandler,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	return cli
}

func messageHandler(topicName, message []byte) {
	fmt.Println(string(topicName), string(message))

	had := &homeAssistantData{}
	err := json.Unmarshal(message, had)
	if err != nil {
		logger.Error("Could not unmarshall data from HomeAssistant: %v", err)
		return
	}
	logger.Infof("type = %s", had.EventType)
	//had.EventData.EntityID name
	/*valueFloat, ok := had.EventData.NewState.Attributes["Temperature"].(float64)
	if !ok {
		logger.Error("Could not parse float", zap.Any("value", had.EventData.NewState.Attributes["Temperature"]))
		return
	}*/

	if had != nil && had.EventData.NewState.EntityID == "sensor.ute_tvistevagen_temperature" {
		temperatureOut = had.EventData.NewState.State
	}
}
