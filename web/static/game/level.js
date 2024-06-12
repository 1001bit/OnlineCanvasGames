class Level {
    constructor(canvas, tps){
        this.canvas = canvas

        setInterval(() => {this.tick()}, 1000/tps)

        this.playersLimit = 0

        this.kinematicRects = new Map([
            [50, canvas.camera]
        ])
        this.staticRects = new Map()

        this.timer = new DeltaTimer()
        this.controls = new Controls()
    }

    setPlayersLimit(limit){
        this.playersLimit = limit
    }

    update(dt){
        // TODO: Something here
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
        this.canvas.deleteDrawable(id)
        this.kinematicRects.delete(id)
        this.staticRects.delete(id)
    }

    tick(){
        let dt = this.timer.getDeltaTime()
        this.update(dt)

        this.canvas.draw()
    }
}