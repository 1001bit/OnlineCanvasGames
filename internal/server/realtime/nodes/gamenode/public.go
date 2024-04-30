package gamenode

// ask gameNode to update gameNode.roomsJSON
func (gameNode *GameNode) RequestUpdatingRoomsJSON() {
	gameNode.roomsJSONUpdateChan <- struct{}{}
}

func (gameNode *GameNode) GetID() int {
	return gameNode.gameID
}
