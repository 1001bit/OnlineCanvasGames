class Platformer {
    private ticker: Ticker;

    private level: Level;

    private canvas: GameCanvas;
    private controls: Controls;
    private websocket: GameWebSocket;

    private serverTPS: number;
    private clientTPS: number;

    constructor(){
        const layers = 2

        this.level = new Level()

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(30, 100, 100))

        this.controls = new Controls()
        this.bindControls()

        this.websocket = new GameWebSocket()
        const gameTitle = $("main").data("game-title")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameTitle, roomID)

        this.serverTPS = 0;
        this.clientTPS = 0;

        this.ticker = new Ticker()
        this.ticker.start(dt => this.tick(dt))
    }

    private bindControls(){
        const controls = this.controls

        controls.bindControl("d", "right")
        controls.bindControl("a", "left")
        controls.bindControl("w", "jump")
        controls.bindControl(" ", "jump")
    }

    private initWebsocket(gameTitle: string, roomID: number){
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

        this.websocket.openConnection(gameTitle, roomID)
    }

    private tick(dt: number) {
        // draw
        this.canvas.draw(this.level.getCameraPosition())

        // level
        this.level.tick(dt, this.controls)
    }

    private stopWithText(text: string){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    private handleLevelMessage(body: LevelMessage){
        this.level.setConfig(body.config)
        this.level.setPlayerRectID(body.playerRectId)
        this.level.setTPS(body.clientTps, body.serverTps)

        this.serverTPS = body.serverTps
        this.clientTPS = body.clientTps

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

    private handleLevelUpdateMessage(body: LevelUpdateMessage){
        this.level.handlePlayerMovement(body.players, body.correct)

        // send controls right after level message, because server allows sending messages right after sending level message
        const heldControlsTicks = this.controls.getHeldControlsTicks()
        if(heldControlsTicks.size != 0){
            // no need of cutting ticks in map, that is being sent to server, since ticks are being limited there
            const json = JSON.stringify(Object.fromEntries(heldControlsTicks))
            this.websocket.sendMessage("input", json)

            // cutting ticks after sending
            this.controls.resetHeldControlsTicks(this.serverTPS, this.clientTPS)
        }
    }

    private handleConnectMessage(body: ConnectMessage){
        let serverRect = body.rect
        let rectID = body.rectId
        
        const rectangle = this.level.createPlayerRectangle(serverRect, rectID)
        if(rectangle){
            this.canvas.insertDrawable(rectangle, 0, rectID)
        }
    }

    private handleDisconnectMessage(body: DisconnectMessage){
        this.canvas.deleteDrawable(body.rectId)
        this.level.disconnectPlayer(body.rectId)
    }
}

new Platformer()