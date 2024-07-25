class Platformer {
    ticker: Ticker;

    level: Level;

    canvas: GameCanvas;
    controls: Controls;
    websocket: GameWebSocket;

    DebugServerDeltaTime: DeltaTimer;

    constructor(){
        const layers = 2

        this.level = new Level()

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(30, 100, 100))

        this.controls = new Controls()
        this.bindControls()

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)

        this.DebugServerDeltaTime = new DeltaTimer();

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
                case "level":
                    this.handleLevelMessage(body);
                    break;

                case "levelUpdate":
                    this.handleLevelUpdateMessage(body);
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
        // level
        this.level.tick(dt, this.controls)

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
        this.level.setTPS(body.clientTps, body.tps)

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

    handleLevelUpdateMessage(body: LevelUpdateMessage){
        // console.log(this.DebugServerDeltaTime.getDeltaTime(), 1000/this.level.serverTPS, body)

        this.level.handlePlayerMovement(body.movedPlayers)

        const heldControlsTicks = this.controls.getHeldControlsTicks()
        if(heldControlsTicks.size != 0){
            const json = JSON.stringify(Object.fromEntries(heldControlsTicks))
            this.websocket.sendMessage("input", json)
            this.controls.clearHeldControlsTicks()
        }
    }

    handleConnectMessage(body: ConnectMessage){
        let serverRect = body.rect
        let rectID = body.rectId
        
        const rectangle = this.level.createPlayerRectangle(serverRect, rectID)
        if(rectangle){
            this.canvas.insertDrawable(rectangle, 0, rectID)
        }
    }

    handleDisconnectMessage(body: DisconnectMessage){
        this.canvas.deleteDrawable(body.rectId)
        this.level.disconnectPlayer(body.rectId)
    }
}

new Platformer()