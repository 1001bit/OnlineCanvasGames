class Platformer {
    constructor(){
        const layers = 2

        this.canSendControls = true

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(30, 100, 100))

        this.controls = new Controls()
        this.bindControls()

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

    createRectangleShape(serverRect, rectID){
        let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y)
        rectangle.setPosition(serverRect.position.x, serverRect.position.y)

        if(this.canvas.drawableExists(rectID)){
            return
        }
        this.canvas.insertDrawable(rectangle, 0, rectID)
    }

    handleLevelMessage(body){
        if(!("kinematic" in body) || !("static" in body)){
            return
        }
    
        let serverKinematicRects = body.kinematic
        let serverStaticRects = body.static
    
        for (const idStr in serverKinematicRects){
            let rectID = parseInt(idStr)
            let serverRect = serverKinematicRects[idStr]
    
            this.createRectangleShape(serverRect, rectID)
        }
        for (const idStr in serverStaticRects){
            let rectID = parseInt(idStr)
            let serverRect = serverStaticRects[idStr]
    
            this.createRectangleShape(serverRect, rectID)
        }
    }

    handleDeleteMessage(body){
        let rectID = parseInt(body)
        this.canvas.deleteDrawable(rectID)
    }

    handleCreateMessage(body){
        let serverRect = body.rect
        let rectID = parseInt(body.id)
        
        this.createRectangleShape(serverRect, rectID, "velocity" in body.rect)
    }

    handleDeltasMessage(body){
        this.canSendControls = true
    
        for (const idStr in body){
            let rectID = parseInt(idStr)
            let serverRect = body[idStr]

            if(!this.canvas.drawableExists(rectID)){
                return
            }
            this.canvas.getDrawable(rectID).setPosition(serverRect.position.x, serverRect.position.y)
        }
    }

    handleGameInfoMessage(body){
        this.playerRectID = body.rectID

        // TODO: Handle constants and make physical calculations based on them
        console.log(body.constants)
    }
}

new Platformer()