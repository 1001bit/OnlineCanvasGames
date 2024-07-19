package gamelogic

import (
	"encoding/json"
)

type UserInput struct {
	Controls map[string]float64
	UserID   int
}

func (input *UserInput) GetControlCoeff(id string) (float64, bool) {
	coeff, ok := input.Controls[id]
	if coeff == 0 || !ok {
		return 0, false
	}

	return min(coeff, 1), true
}

func (input *UserInput) IsHeld(id string) bool {
	coeff, ok := input.Controls[id]
	return coeff != 0 && ok
}

func ExtractInputFromMsg(body any, userID int, inputChan chan<- UserInput) {
	inputMap := make(map[string]float64)

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
