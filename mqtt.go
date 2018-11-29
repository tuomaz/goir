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
	cli := client.New(&client.Options{
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	logger.Infof("1")

	err := cli.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  serverAddress,
		ClientID: []byte(clientID),
	})
	if err != nil {
		panic(err)
	}

	logger.Infof("3")
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(topic),
				QoS:         mqtt.QoS0,
				Handler:     messageHandler,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	logger.Infof("3")
	return cli
}

func messageHandler(topicName, message []byte) {
	had := &homeAssistantData{}
	err := json.Unmarshal(message, had)
	if err != nil {
		logger.Error("Could not unmarshall data from HomeAssistant: %v", err)
		return
	}
	logger.Debug("data = %v", had)
	//had.EventData.EntityID name
	/*valueFloat, ok := had.EventData.NewState.Attributes["Temperature"].(float64)
	if !ok {
		logger.Error("Could not parse float", zap.Any("value", had.EventData.NewState.Attributes["Temperature"]))
		return
	}*/
}
