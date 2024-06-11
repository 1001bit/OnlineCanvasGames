const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")
const layers = 2

const game = new Game(gameID, roomID, layers)
game.canvas.setBackgroundColor(RGB(30, 200, 200))

const level = new Level(game.canvas)

level.controls.bindControl("d", "right")
level.controls.bindControl("a", "left")
level.controls.bindControl("w", "up")
level.controls.bindControl("s", "down")

let playerRectID = -1

function handleLevelMessage(body){
    if(!("kinematic" in body) || !("static" in body)){
        return
    }

    level.updateKinematics()
    game.websocket.sendMessage("input", level.controls.getControlsJSON())
    level.controls.clear()

    let kinematic = body.kinematic
    let static = body.static

    for (idStr in kinematic){
        let rectID = parseInt(idStr)
        let clientRect = 0;
        let serverRect = kinematic[idStr]

        if (!level.kinematicRects.has(rectID)){
            let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y, true)
            level.insertDrawable(rectangle, 0, rectID)
            clientRect = rectangle.rect
        } else {
            clientRect = level.kinematicRects.get(rectID)
        }
        
        clientRect.setTargetPos(serverRect.position.x, serverRect.position.y, false)
    }
    for (idStr in static){
        let rectID = parseInt(idStr)
        let clientRect = 0;
        let serverRect = static[idStr]

        if (!level.kinematicRects.has(rectID)){
            let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y, false)
            level.insertDrawable(rectangle, 0, rectID)
            clientRect = rectangle.rect
        } else {
            clientRect = level.kinematicRects.get(rectID)
        }

        clientRect.setPosition(serverRect.position.x, serverRect.position.y)
    }
}

function handleGameInfoMessage(body){
    level.setTPS(60, body.tps)
    level.setPlayersLimit(body.limit)
    playerRectID = body.rectID
}

function handleDeleteMessage(body){
    level.deleteDrawable(parseInt(body))
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

        case "delete":
            handleDeleteMessage(body);
            break;

        default:
            break;
    }
}