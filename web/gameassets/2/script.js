class Platformer {
    constructor(){
        const layers = 2

        this.game = new Game(layers)
        this.game.canvas.setBackgroundColor(RGB(30, 200, 200))

        this.controls = new Controls()
        this.bindControls()

        this.level = new Level(this.game.canvas)

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)
    }

    bindControls(){
        const controls = this.controls

        controls.bindControl("d", "right")
        controls.bindControl("a", "left")
        controls.bindControl("w", "jump")
        controls.bindControl(" ", "jump")
    }

    initWebsocket(gameID, roomID){
        this.game.initWebsocket(this.websocket, gameID, roomID, (type, body) => {
            switch (type) {
                case "gameinfo":
                    this.handleGameInfoMessage(body);
                    break;
        
                case "level":
                    this.handleLevelMessage(body);
                    break;
        
                case "deltas":
                    this.handleDeltasMessage(body);
                    break;
        
                case "delete":
                    this.handleDeleteMessage(body);
                    break;
        
                case "create":
                    this.handleCreateMessage(body);
                    break;
        
                default:
                    break;
            }
        })
    }

    createRect(serverRect, rectID, kinematic){
        const level = this.level

        let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y, kinematic)
        rectangle.rect.setPosition(serverRect.position.x, serverRect.position.y)
        level.insertDrawable(rectangle, 0, rectID)
    }

    handleLevelMessage(body){
        if(!("kinematic" in body) || !("static" in body)){
            return
        }
    
        let kinematicRects = body.kinematic
        let staticRects = body.static
    
        for (const idStr in kinematicRects){
            let rectID = parseInt(idStr)
            let serverRect = kinematicRects[idStr]
    
            this.createRect(serverRect, rectID, true)
        }
        for (const idStr in staticRects){
            let rectID = parseInt(idStr)
            let serverRect = staticRects[idStr]
    
            this.createRect(serverRect, rectID, false)
        }
    }

    handleDeleteMessage(body){
        this.level.deleteDrawable(parseInt(body))
    }

    handleCreateMessage(body){
        const level = this.level

        let serverRect = body.rect
        let rectID = parseInt(body.id)
    
        if (level.kinematicRects.has(rectID) || level.staticRects.has(rectID)){
            return
        }
    
        if ("velocity" in body.rect){
            this.createRect(serverRect, rectID, true)
        } else {
            this.createRect(serverRect, rectID, false)
        }
    }

    handleDeltasMessage(body){
        const level = this.level

        this.websocket.sendMessage("input", this.controls.getControlsJSON())
    
        for (const idStr in body){
            let rectID = parseInt(idStr)
            if(!level.kinematicRects.has(rectID)){
                continue
            }
            let serverRect = body[idStr]
    
            level.kinematicRects.get(rectID).setPosition(serverRect.position.x, serverRect.position.y)
        }
    }

    handleGameInfoMessage(body){
        this.playerRectID = body.rectID
    }
}

const platformer = new Platformer()