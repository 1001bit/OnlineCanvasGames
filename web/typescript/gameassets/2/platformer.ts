/// <reference path="physicsEngine.ts"/>

class Platformer {
    playerRectID: number;

    controlsAccumulator: number;
    serverTPS: number;

    canvas: GameCanvas;
    controls: Controls;
    websocket: GameWebSocket;
    ticker: Ticker;

    physicsEngine: Physics;

    constructor(){
        const layers = 2

        this.playerRectID = 0

        this.controlsAccumulator = 0;
        this.serverTPS = 0

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(30, 100, 100))

        this.controls = new Controls()
        this.bindControls()

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)
        
        this.physicsEngine = new Physics()

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

    tick(dt: number) {
        this.physicsEngine.tick(dt)

        this.controlsAccumulator += dt
        const maxControlsAccumulator = 1000/(this.serverTPS*4)
        while(this.controlsAccumulator > maxControlsAccumulator){
            let heldControls = this.controls.getHeldControls()
            if(heldControls.size > 0){
                let json = JSON.stringify(Object.fromEntries(heldControls.entries()))
                this.websocket.sendMessage("input", json)
            }
            this.controlsAccumulator -= maxControlsAccumulator
        }

        this.canvas.draw()
    }

    stopWithText(text: string){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    createRectangleShape(serverRect: Rect | KinematicRect, rectID: number){
        if(this.canvas.getDrawable(rectID)){
            return
        }

        let rectangle: RectangleShape;

        if(isKinematicRect(serverRect)){
            const rect = new KinematicRect(serverRect)
            rect.setVelocity(serverRect.velocity.x, serverRect.velocity.y)
            this.physicsEngine.insertKinematicRect(rectID, rect)

            rectangle = new RectangleShape(rect)
        } else {
            const rect = new Rect(serverRect)
            this.physicsEngine.insertStaticRect(rectID, rect)

            rectangle = new RectangleShape(rect)
        }

        this.canvas.insertDrawable(rectangle, 0, rectID)
    }

    handleLevelMessage(body: LevelMessage){
        for (const [key, val] of Object.entries(body.kinematic)){
            const id = Number(key)
            const serverRect = val as KinematicRect

            this.createRectangleShape(serverRect, id)
        }

        for (const [key, val] of Object.entries(body.static)){
            const id = Number(key)
            const serverRect = val as Rect

            this.createRectangleShape(serverRect, id)
        }
    }

    handleDeleteMessage(body: DeleteMessage){
        this.canvas.deleteDrawable(body.ID)
        this.physicsEngine.deleteRect(body.ID)
    }

    handleCreateMessage(body: CreateMessage){
        let serverRect = body.rect
        let rectID = body.id
        
        this.createRectangleShape(serverRect, rectID)
    }

    handleGameInfoMessage(body: GameInfoMessage){
        this.playerRectID = body.rectID
        this.serverTPS = body.tps

        this.physicsEngine.setConstants(body.constants)
    }
}

new Platformer()