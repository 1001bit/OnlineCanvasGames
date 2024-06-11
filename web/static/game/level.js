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
        this.accumulator = 0
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
        this.canvas.insertDrawable(drawable, layer, id)

        if(drawable.rect.isKinematic()){
            this.kinematicRects.set(id, drawable.rect)
        } else {
            this.staticRects.set(id, drawable.rect)
        }
    }

    deleteDrawable(id){
        this.canvas.drawablesLayers.deleteDrawable(id)
        this.kinematicRects.delete(id)
        this.staticRects.delete(id)
    }

    tick(){
        let dt = this.timer.getDeltaTime()
        this.update(dt)

        this.canvas.draw()
    }
}