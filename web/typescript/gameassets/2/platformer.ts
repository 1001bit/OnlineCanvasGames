class Platformer {
    serverTPS: number;

    ticker: Ticker;

    level: Level;

    canvas: GameCanvas;
    controls: Controls;
    websocket: GameWebSocket;

    constructor(){
        const layers = 2

        this.serverTPS = 0;

        this.level = new Level()

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

                case "playerMovement":
                    this.handlePlayerMovementMessage(body);
                    break;
                
                case "connect":
                    this.handleConnectMessage(body);
                    break;

                case "disconnect":
                    this.handleDisconnectMessage(body);
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

    tick(dt: number) {
        // controls
        this.controls.updateCoeffs(this.serverTPS, 1000/dt)

        // level
        this.level.tick(dt, this.serverTPS, this.controls)

        // draw
        this.canvas.draw()
    }

    stopWithText(text: string){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    handleLevelMessage(body: LevelMessage){
        this.level.setConfig(body.config)
        this.level.setPlayerRectID(body.playerRectId)

        for (const [key, val] of Object.entries(body.players)){
            const id = Number(key)
            const serverRect = val as AbstractPlayer

            const rectangle = this.level.createPlayerRectangle(serverRect, id)
            if (rectangle){
                this.canvas.insertDrawable(rectangle, 0, id)
            }
        }

        for (const [key, val] of Object.entries(body.blocks)){
            const id = Number(key)
            const serverRect = val as AbstractBlock

            const rectangle = this.level.createBlockRectangle(serverRect, id)
            if (rectangle){
                this.canvas.insertDrawable(rectangle, 0, id)
            }
        }
    }

    handlePlayerMovementMessage(body: PlayerMovementMessage){
        // level
        this.level.handlePlayerMovement(body.playersMoved)

        // send controls to server
        const controlsCoeffs = this.controls.getCoeffs()
        if(controlsCoeffs.size > 0){
            const json = JSON.stringify(Object.fromEntries(controlsCoeffs.entries()))
            this.controls.resetCoeffs()
            this.websocket.sendMessage("input", json)
        }
    }

    handleDisconnectMessage(body: DisconnectMessage){
        this.canvas.deleteDrawable(body.rectId)
        this.level.disconnectPlayer(body.rectId)
    }

    handleConnectMessage(body: ConnectMessage){
        let serverRect = body.rect
        let rectID = body.rectId
        
        const rectangle = this.level.createPlayerRectangle(serverRect, rectID)
        if(rectangle){
            this.canvas.insertDrawable(rectangle, 0, rectID)
        }
    }

    handleGameInfoMessage(body: GameInfoMessage){
        this.serverTPS = body.tps
    }
}

new Platformer()