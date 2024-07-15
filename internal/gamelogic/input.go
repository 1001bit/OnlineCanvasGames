package gamelogic

import (
	"encoding/json"
)

type UserInput struct {
	Controls map[string]bool
	UserID   int
}

func (input UserInput) IsControlHeld(id string) bool {
	_, ok := input.Controls[id]
	return ok
}

func ExtractInputFromMsg(body any, userID int, inputChan chan<- UserInput) {
	inputMap := make(map[string]bool)

	err := json.Unmarshal([]byte(body.(string)), &inputMap)
	if err != nil {
		return
	}

	go func() {
		inputChan <- UserInput{
			Controls: inputMap,
			UserID:   userID,
		}
	}()
}
