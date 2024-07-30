package platformer

import "github.com/neinBit/ocg-games-service/internal/gamelogic"

type PlayerData struct {
	player *Player
	rectID rectID

	input *PlayerInput
}

func NewPlayerData(player *Player, rectID rectID) *PlayerData {
	return &PlayerData{
		player: player,
		rectID: rectID,

		input: nil,
	}
}

func (pd *PlayerData) HandleInput(inputMap gamelogic.InputMap, serverTPS, clientTPS float64) {
	pd.input = NewPlayerInput(inputMap, serverTPS, clientTPS)
}

func (pd *PlayerData) ControlPlayer(speed, jump float64) {
	pd.player.Control(speed, jump, pd.input)
	pd.input = nil
}
