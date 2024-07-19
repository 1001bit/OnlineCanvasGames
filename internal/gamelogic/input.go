package gamelogic

import (
	"encoding/json"
)

type UserInput struct {
	inputMap map[string]float64
	userID   int
}

func NewUserInput(userID int, inputMap map[string]float64) *UserInput {
	return &UserInput{
		userID:   userID,
		inputMap: inputMap,
	}
}

func (input *UserInput) GetControlCoeff(id string) (float64, bool) {
	coeff, ok := input.inputMap[id]
	if coeff == 0 || !ok {
		return 0, false
	}

	return min(coeff, 1), true
}

func (input *UserInput) IsHeld(id string) bool {
	coeff, ok := input.inputMap[id]
	return coeff != 0 && ok
}

func (input *UserInput) GetUserID() int {
	return input.userID
}

func GetInputMapFromMsg(body any, userID int) (map[string]float64, error) {
	inputMap := make(map[string]float64)

	err := json.Unmarshal([]byte(body.(string)), &inputMap)
	if err != nil {
		return nil, err
	}

	return inputMap, nil
}
