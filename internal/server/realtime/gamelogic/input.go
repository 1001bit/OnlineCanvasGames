package gamelogic

import (
	"encoding/json"
)

type Control struct {
	HoldPeriod float64 `json:"holdPeriod"`
}

type UserInput struct {
	Controls map[string]Control
	UserID   int
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
