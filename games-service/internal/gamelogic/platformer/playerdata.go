package platformer

import "github.com/neinBit/ocg-games-service/internal/gamelogic"

type PlayerData struct {
	player *Player
	rectID rectID

	input PlayerInput
}

func NewPlayerData(player *Player, rectID rectID, serverTPS, clientTPS float64) *PlayerData {
	return &PlayerData{
		player: player,
		rectID: rectID,

		input: CreatePlayerInput(serverTPS, clientTPS),
	}
}

func (pd *PlayerData) HandleInput(inputMap gamelogic.InputMap) {
	pd.input.SetInputMap(inputMap)
}

func (pd *PlayerData) ControlPlayer(speed, jump float64) {
	pd.player.Control(speed, jump, &pd.input)
	pd.input.ClearInputMap()
}
