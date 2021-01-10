package msg

import (
	"chanv/link-monitor/cmd"
	"encoding/json"
	"fmt"
	"time"
)

func StartPublishing(msgClient connection, routingKey string, publish_interval uint, links []*cmd.Iface){
	for {
		msg := getPublishMessage(links)

		fmt.Println(msg)
		msgClient.Publish(routingKey, msg)

		time.Sleep(time.Duration(publish_interval) * time.Second)
	}
}

func getPublishMessage(links []*cmd.Iface) string {
	message := map[string] cmd.StatMessage {}

	for _, link := range links {
		message[link.Name()] = link.GetIfaceStatsMessage()
	}

	messageJson, _ := json.Marshal(message)

	return string(messageJson)
}
