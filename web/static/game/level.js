class Level {
    constructor(canvas, drawables){
        this.drawables = drawables
        this.canvas = canvas

        this.serverTPS = 0
        this.clientTPS = 0
        this.accumulator = 0

        this.timer = new DeltaTimer()
        this.controls = new Controls()
    }

    handleLevelMessage(body, websocket) {
        for (const key in body){
            if(!(key in this.drawables)){
                continue
            }
            
            this.drawables[key].setPosition(body[key].x, body[key].y)
            this.drawables[key].setSize(body[key].w, body[key].h)
        }

        websocket.sendMessage("input", this.controls.getControlsJSON())
        this.controls.clear()
    }

    setTPS(client, server){
        this.clientTPS = client
        this.serverTPS = server

        this.timer.getDeltaTime()
        setInterval(() => {this.tick()}, 1000/this.clientTPS)
    }

    update(dt){
        this.controls.updateHoldTime(dt)
        
        this.accumulator += dt
        let alpha = this.accumulator/(1000/this.serverTPS)
        // TODO: Interpolate all the rects, that have been transformed

        this.controls.updatePressedStatus()
    }

    tick(){
        let dt = this.timer.getDeltaTime()
        this.update(dt)

        this.canvas.draw()
    }
}