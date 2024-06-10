class Level {
    constructor(canvas){
        this.canvas = canvas

        this.serverTPS = 0
        this.clientTPS = 0
        this.accumulator = 0

        this.playersLimit = 0

        this.kinematicRects = new Map([
            [50, canvas.camera]
        ])
        this.staticRects = new Map()

        this.timer = new DeltaTimer()
        this.controls = new Controls()
    }

    handleLevelMessage(body, websocket) {
        this.updateKinematics()
        this.accumulator = 0

        for (const key in body){
            let rectID = parseInt(key)
            if(this.staticRects.has(rectID)){
                this.staticRects.get(rectID).setPosition(body[key].x, body[key].y)
                this.staticRects.get(rectID).setSize(body[key].w, body[key].h)
            }
            if(this.kinematicRects.has(rectID)){
                this.kinematicRects.get(rectID).setTargetPos(body[key].x, body[key].y)
                this.kinematicRects.get(rectID).setSize(body[key].w, body[key].h)
            }
        }

        websocket.sendMessage("input", this.controls.getControlsJSON())
        this.controls.clear()

        return body
    }

    setTPS(client, server){
        this.clientTPS = client
        this.serverTPS = server

        this.timer.getDeltaTime()
        setInterval(() => {this.tick()}, 1000/this.clientTPS)
    }

    setPlayersLimit(limit){
        this.playersLimit = limit
    }

    update(dt){
        this.controls.updateHoldTime(dt)
        this.accumulator += dt
        let alpha = this.accumulator/(1000/this.serverTPS)
        this.interpolateKinematic(alpha)

        this.controls.updatePressedStatus()
    }

    updateKinematics(){
        this.kinematicRects.forEach((kinRect) => {
            kinRect.updatePrevPos()
        })
    }

    interpolateKinematic(alpha){
        alpha = Math.min(alpha, 1)
        this.kinematicRects.forEach((kinRect) => {
            kinRect.interpolate(alpha)
        })
    }

    insertDrawable(drawable, layer, id){
        this.canvas.drawablesLayers.insertDrawable(drawable, layer)

        if(drawable.rect.isKinematic()){
            this.kinematicRects.set(id, drawable.rect)
        } else {
            this.staticRects.set(id, drawable.rect)
        }
    }

    rectExists(id){
        return this.kinematicRects.has(id) || this.staticRects.has(id)
    }

    tick(){
        let dt = this.timer.getDeltaTime()
        this.update(dt)

        this.canvas.draw()
    }
}