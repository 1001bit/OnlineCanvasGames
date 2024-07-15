class Platformer {
    canvas: GameCanvas;
    controls: Controls;
    websocket: GameWebSocket;
    ticker: Ticker;

    canSendControls: boolean;
    playerRectID: number;

    constructor(){
        const layers = 2

        this.playerRectID = 0
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

    initWebsocket(gameID: number, roomID: number){
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

    tick(_dt: number) {
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

    stopWithText(text: string){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    createRectangleShape(serverRect: Rect, rectID: number){
        let rectangle = new RectangleShape(serverRect.size.x, serverRect.size.y)
        rectangle.setPosition(serverRect.position.x, serverRect.position.y)

        this.canvas.insertDrawable(rectangle, 0, rectID)
    }

    handleLevelMessage(body: LevelMessage){
        for (let [key, val] of Object.entries(body.kinematic)){
            const id = Number(key)
            const serverRect = val as Rect

            this.createRectangleShape(serverRect, id)
        }

        for (let [key, val] of Object.entries(body.static)){
            const id = Number(key)
            const serverRect = val as Rect

            this.createRectangleShape(serverRect, id)
        }
    }

    handleDeleteMessage(body: DeleteMessage){
        this.canvas.deleteDrawable(body.ID)
    }

    handleCreateMessage(body: CreateMessage){
        let serverRect = body.rect
        let rectID = body.id
        
        this.createRectangleShape(serverRect, rectID)
    }

    handleDeltasMessage(body: DeltasMessage){
        this.canSendControls = true

        for (let [key, val] of Object.entries(body.kinematic)){
            const id = Number(key)
            const serverRect = val as Rect

            const clientRect = this.canvas.getDrawable(id)

            if(clientRect){
                clientRect.setPosition(serverRect.position.x, serverRect.position.y)
            }
        }
    }

    handleGameInfoMessage(body: GameInfoMessage){
        this.playerRectID = body.rectID

        // TODO: Handle constants and make physical calculations based on them
        console.log(body.constants)
    }
}

new Platformer()