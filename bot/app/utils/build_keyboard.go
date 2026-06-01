package utils

import (
	"encoding/json"
	"log"
)

func MakeReplyMarkup(jsonStr string) json.RawMessage {
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		log.Println("Error in MakeReplyMarkup:", err)
		return nil
	}
	return raw
}
