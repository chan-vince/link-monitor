package cmd

import (
	"strings"
)

func ConstructRoutingKey(baseKey string, kitId string) string {
	return strings.Replace(baseKey, "<kit_id>", kitId, -1)
}