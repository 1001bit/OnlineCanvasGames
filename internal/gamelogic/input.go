package gamelogic

import "encoding/json"

type InputMap map[string]float64

func GetInputMapFromMsg(body any, userID int) (InputMap, error) {
	inputMap := make(InputMap)

	err := json.Unmarshal([]byte(body.(string)), &inputMap)
	if err != nil {
		return nil, err
	}

	return inputMap, nil
}

func (inputMap InputMap) GetControlCoeff(id string) (float64, bool) {
	coeff, ok := inputMap[id]
	if coeff == 0 || !ok {
		return 0, false
	}

	return min(coeff, 1), true
}

func (inputMap InputMap) IsHeld(id string) bool {
	coeff, ok := inputMap[id]
	return coeff != 0 && ok
}
