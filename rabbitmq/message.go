package rabbitmq

import (
	"chanv/link-monitor/cmd"
	"encoding/json"
	"time"
)

func StartPublishing(msgClient *client, routingKey string, publish_interval uint, links []*cmd.Iface){
	for {
		msg := GetPublishMessage(links)

		//log.Println(msg)
		msgClient.Publish(routingKey, msg)

		time.Sleep(time.Duration(publish_interval) * time.Second)
	}
}

func GetPublishMessage(links []*cmd.Iface) string {
	message := map[string] cmd.StatMessage {}

	for _, link := range links {
		message[link.Name()] = link.GetIfaceStatsMessage()
	}

	messageJson, _ := json.Marshal(message)

	return string(messageJson)
}
