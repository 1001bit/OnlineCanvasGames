package platformer

import (
	"encoding/json"
)

type Control struct {
	WasPressed bool `json:"wasPressed"`
	HoldPeriod int  `json:"holdPeriod"`
}

func (gl *PlatformerGL) handleInputMessage(body any, userID int) {
	inputMap := make(map[string]Control)

	err := json.Unmarshal([]byte(body.(string)), &inputMap)
	if err != nil {
		return
	}

	_ = userID
	// TODO: Use input for something
}
