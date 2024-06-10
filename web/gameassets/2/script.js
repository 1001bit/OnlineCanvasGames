const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")

const game = new Game(gameID, roomID)
game.canvas.setLayersCount(2)
game.canvas.setBackgroundColor(RGB(30, 200, 200))

const level = new Level(game.canvas)

level.controls.bindControl("d", "right")
level.controls.bindControl("a", "left")
level.controls.bindControl("w", "up")
level.controls.bindControl("s", "down")

let playerRectID = -1

function handleLevelMessage(body){
    // get all the rects from the message
    let rects = level.handleLevelMessage(body, game.websocket)

    // iterate over all the rects
    for (const key in rects){
        let rectID = parseInt(key)

        // is rectID is in "player" range, and the rect doesn't exist yet, create new player rect
        if (rectID < level.playersLimit && !level.rectExists(rectID)){
            const player = new RectangleShape(0, 0, true)
            level.insertDrawable(player, 0, rectID)
        }
    }

    // bind camera to player rect if it exists
    if (level.rectExists(playerRectID)){
        let playerPos = level.kinematicRects.get(playerRectID).getPosition()
        // game.canvas.setCameraPos(playerPos.x - this.canvas.width/2 + 50, playerPos.y - this.canvas.height/2 + 50)
    }
}

function handleGameInfoMessage(body){
    level.setTPS(90, body.tps)
    level.setPlayersLimit(body.limit)
    playerRectID = body.rectID
}

// on server message
game.handleGameMessage = (type, body) => {
    switch (type) {
        case "gameinfo":
            handleGameInfoMessage(body);
            break;

        case "level":
            handleLevelMessage(body);
            break;

        default:
            break;
    }
}