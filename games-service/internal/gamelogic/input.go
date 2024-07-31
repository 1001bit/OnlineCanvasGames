package gamelogic

// inputMap[control]ticks
type InputMap map[string]int

func (inputMap InputMap) IsHeld(control string) bool {
	_, ok := inputMap[control]
	return ok
}

func (inputMap InputMap) GetTicks(control string) (int, bool) {
	ticks, ok := inputMap[control]
	return ticks, ok
}
