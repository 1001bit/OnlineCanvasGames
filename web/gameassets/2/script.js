const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")

const game = new Game(gameID, roomID)
game.canvas.setLayersCount(2)
game.canvas.setBackgroundColor(RGB(30, 200, 200))

const rectangle = new RectangleShape(0, 0, false)
game.canvas.drawablesLayers.insertDrawable(rectangle, 0)
const level = new Level({1: rectangle})

// on server message
game.handleGameMessage = (type, body) => {
    // TODO: Handle gameinfo type

    if(type == "level"){
        level.handleLevelMessage(body)
    }
}