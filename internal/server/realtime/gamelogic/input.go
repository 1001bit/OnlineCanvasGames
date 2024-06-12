package gamelogic

import (
	"encoding/json"
)

type Control struct {
	IsHeld bool `json:"isHeld"`
}

type UserInput struct {
	Controls map[string]Control
	UserID   int
}

func (input UserInput) IsControlHeld(id string) bool {
	control, ok := input.Controls[id]
	return ok && control.IsHeld
}

func ExtractInputFromMsg(body any, userID int, inputChan chan<- UserInput) {
	inputMap := make(map[string]Control)

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
