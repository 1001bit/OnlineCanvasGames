const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")
const layers = 2

const game = new Game(gameID, roomID, layers)
game.canvas.setBackgroundColor(RGB(30, 200, 200))

const level = new Level(game.canvas)

level.controls.bindControl("d", "right")
level.controls.bindControl("a", "left")
level.controls.bindControl("w", "jump")
level.controls.bindControl(" ", "jump")

let playerRectID = -1

function createRect(serverRect, rectID, kinematic){
    let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y, kinematic)
    rectangle.rect.setPosition(serverRect.position.x, serverRect.position.y)
    level.insertDrawable(rectangle, 0, rectID)
}

function handleLevelMessage(body){
    if(!("kinematic" in body) || !("static" in body)){
        return
    }

    let kinematic = body.kinematic
    let static = body.static

    for (idStr in kinematic){
        let rectID = parseInt(idStr)
        let serverRect = kinematic[idStr]

        createRect(serverRect, rectID, true)
    }
    for (idStr in static){
        let rectID = parseInt(idStr)
        let serverRect = static[idStr]

        createRect(serverRect, rectID, false)
    }
}

function handleDeleteMessage(body){
    level.deleteDrawable(parseInt(body))
}

function handleCreateMessage(body){
    let serverRect = body.rect
    let rectID = parseInt(body.id)

    if (level.kinematicRects.has(rectID) || level.staticRects.has(rectID)){
        return
    }

    if ("velocity" in body.rect){
        createRect(serverRect, rectID, true)
    } else {
        createRect(serverRect, rectID, false)
    }
}

function handleDeltasMessage(body){
    game.websocket.sendMessage("input", level.controls.getControlsJSON())

    for (idStr in body){
        let rectID = parseInt(idStr)
        if(!level.kinematicRects.has(rectID)){
            continue
        }
        let serverRect = body[idStr]

        level.kinematicRects.get(rectID).setPosition(serverRect.position.x, serverRect.position.y)
    }
}

level.updateKinematics = dt => {
    
}

function handleGameInfoMessage(body){
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

        case "deltas":
            handleDeltasMessage(body);
            break;

        case "delete":
            handleDeleteMessage(body);
            break;

        case "create":
            handleCreateMessage(body);
            break;

        default:
            break;
    }
}