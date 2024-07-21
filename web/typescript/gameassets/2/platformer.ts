class Platformer {
    playerRectID: number;
    constants: PlatformerConstants;

    physTps: number;
    physTicker: FixedTicker;

    serverAccumulator: number;
    serverTPS: number;

    canvas: GameCanvas;
    controls: Controls;
    websocket: GameWebSocket;
    ticker: Ticker;

    physicsEngine: PhysicsEngine;

    constructor(){
        const layers = 2

        this.playerRectID = 0
        this.constants = {
            physics: {
                friction: 0,
                gravity: 0,
            },
            playerSpeed: 0,
            playerJump: 0,
        }

        this.physTps = 30;
        this.physTicker = new FixedTicker(this.physTps);

        this.serverAccumulator = 0;
        this.serverTPS = 0;

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(30, 100, 100))

        this.controls = new Controls()
        this.bindControls()

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)
        
        this.physicsEngine = new PhysicsEngine()

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

                case "update":
                    this.handleUpdateMessage(body);
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
        // physics
        this.physTicker.update(dt, (fixedDT) => {
            // update phys/server tps coeffs
            this.controls.updateCoeffs(this.serverTPS, this.physTps)

            this.physicsEngine.updateKinematicsInterpolation()
            this.handleControls()
            this.physicsEngine.update(fixedDT, this.constants.physics)
        })

        // interpolation
        this.serverAccumulator += dt
        const interpolatedAlpha = Math.min(1, this.serverAccumulator/(1000/this.serverTPS))
        const kinematicAlpha = this.physTicker.getAlpha()
        this.physicsEngine.interpolate(interpolatedAlpha, kinematicAlpha)

        // draw
        this.canvas.draw()
    }

    handleControls(){
        const playerRect = this.physicsEngine.kinematicRects.get(this.playerRectID)
        if(!playerRect){
            return
        }

        if(this.controls.isHeld("left")){
            playerRect.velocity.x -= this.constants.playerSpeed
        }
        if(this.controls.isHeld("right")){
            playerRect.velocity.x += this.constants.playerSpeed
        }
        if(this.controls.isHeld("jump") && playerRect.collisionVertical == Direction.Down){
            playerRect.velocity.y -= this.constants.playerJump
        }
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

        if(isAbstractKinematicRect(serverRect) && rectID == this.playerRectID){
            // Doing physics prediction only for player rect
            const rect = new KinematicRect(serverRect)
            this.physicsEngine.insertKinematicRect(rectID, rect)

            rectangle = new RectangleShape(rect)
        } else if(isAbstractKinematicRect(serverRect)) {
            // Interpolated rect for other kinematic rects
            const rect = new InterpolatedRect(serverRect)
            this.physicsEngine.insertInterpolatedRect(rectID, rect)

            rectangle = new RectangleShape(rect)
        } else if(isAbstractPhysicalRect(serverRect)) {
            // Static rects
            const rect = new PhysicalRect(serverRect)
            this.physicsEngine.insertStaticRect(rectID, rect)

            rectangle = new RectangleShape(rect)
        } else {
            // non rects
            return
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

    handleUpdateMessage(body: UpdateMessage){
        // physics operations
        this.serverAccumulator = 0
        this.physicsEngine.updateInterpolatedInterpolation()
        this.physicsEngine.setMultiplePositions(body.rectsMoved)

        // send controls to server
        const controlsCoeffs = this.controls.getCoeffs()
        if(controlsCoeffs.size > 0){
            const json = JSON.stringify(Object.fromEntries(controlsCoeffs.entries()))
            this.controls.resetCoeffs()
            this.websocket.sendMessage("input", json)
        }
    }

    handleDeleteMessage(body: DeleteMessage){
        this.canvas.deleteDrawable(body.id)
        this.physicsEngine.deleteRect(body.id)
    }

    handleCreateMessage(body: CreateMessage){
        let serverRect = body.rect
        let rectID = body.id
        
        this.createRectangleShape(serverRect, rectID)
    }

    handleGameInfoMessage(body: GameInfoMessage){
        this.playerRectID = body.rectID
        this.serverTPS = body.tps
        this.constants = body.constants
    }
}

new Platformer()