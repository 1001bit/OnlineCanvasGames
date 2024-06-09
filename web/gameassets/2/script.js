const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")

const game = new Game(gameID, roomID)
game.canvas.setLayersCount(2)
game.canvas.setBackgroundColor(RGB(30, 200, 200))

const rectangle = new RectangleShape(0, 0, false)
game.canvas.drawablesLayers.insertDrawable(rectangle, 0)
const level = new Level(game.canvas, {1: rectangle})

level.controls.bindControl("d", "right")
level.controls.bindControl("a", "left")

// on server message
game.handleGameMessage = (type, body) => {
    switch (type) {
        case "gameinfo":
            level.setTPS(60, body.tps)
            break;

        case "level":
            level.handleLevelMessage(body, game.websocket)
            break

        default:
            break;
    }
}