package platformer

import "github.com/1001bit/OnlineCanvasGames/internal/gamelogic"

type PlayerData struct {
	player *Player
	rectID int

	input *PlayerInput
}

func NewPlayerData(player *Player, rectID int) *PlayerData {
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
