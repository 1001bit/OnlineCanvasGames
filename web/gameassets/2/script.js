class Platformer {
    constructor(){
        const layers = 2

        this.canSendControls = true

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(30, 100, 100))

        this.controls = new Controls()
        this.bindControls()

        this.level = new Level()

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)

        this.ticker = new Ticker()
        this.ticker.tick(dt => this.tick(dt))
    }

    bindControls(){
        const controls = this.controls

        controls.bindControl("d", "right")
        controls.bindControl("a", "left")
        controls.bindControl("w", "jump")
        controls.bindControl(" ", "jump")
    }

    initWebsocket(gameID, roomID){
        this.websocket.handleMessage = (type, body) => {
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
        }

        this.websocket.handleClose = (body) => {
            this.stopWithText(body)
        }

        this.websocket.openConnection(gameID, roomID)
    }

    tick(dt) {
        if(this.canSendControls){
            let heldControls = this.controls.getHeldControls()
            if(heldControls.size > 0){
                let json = JSON.stringify(Object.fromEntries(heldControls.entries()))
                this.websocket.sendMessage("input", json)
                this.canSendControls = false
            }
        }

        this.canvas.draw()
    }

    stopWithText(text){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    createRectangleShape(serverRect, rectID, kinematic){
        const level = this.level

        let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y, kinematic)
        rectangle.setPosition(serverRect.position.x, serverRect.position.y)

        if(kinematic){
            level.insertKinematicRect(rectID, rectangle.rect)
        } else {
            level.insertStaticRect(rectID, rectangle.rect)
        }

        this.canvas.insertDrawable(rectangle, 0, rectID)
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
    
            this.createRectangleShape(serverRect, rectID, true)
        }
        for (const idStr in staticRects){
            let rectID = parseInt(idStr)
            let serverRect = staticRects[idStr]
    
            this.createRectangleShape(serverRect, rectID, false)
        }
    }

    handleDeleteMessage(body){
        let rectID = parseInt(body)

        this.level.deleteRect(rectID)
        this.canvas.deleteDrawable(rectID)
    }

    handleCreateMessage(body){
        const level = this.level

        let serverRect = body.rect
        let rectID = parseInt(body.id)
    
        if (level.kinematicRects.has(rectID) || level.staticRects.has(rectID)){
            return
        }
        
        this.createRectangleShape(serverRect, rectID, "velocity" in body.rect)
    }

    handleDeltasMessage(body){
        const level = this.level

        this.canSendControls = true
    
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

new Platformer()