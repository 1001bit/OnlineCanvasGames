class Level{
    constructor(layers){
        this.drawablesLayers = new DrawableLayers(layers)
        this.camera = new KinematicRect()
        this.kinematicRects = [this.camera]

        this.updateRate = 20

        this.accumulator = 0
        this.tickRate = 60
        this.tickInterval = setInterval(() => this.tick(), 1000/this.tickRate)
        this.timer = new DeltaTimer()
    }

    tick(){
        let dt = this.timer.getDeltaTime()

        this.update(dt)
    }

    draw(ctx){
        ctx.save()
        ctx.translate(-this.camera.left, -this.camera.top) // for some reason, it has to be a negatile value

        this.drawablesLayers.draw(ctx)

        ctx.restore()
    }

    update(dt){
        this.accumulator += dt

        while(this.accumulator >= 1000/this.updateRate){
            this.accumulator -= 1000/this.updateRate

            this.kinematicRects.forEach(rect => {
                rect.updatePrevPos()
            })
        }

        this.interpolateKinematics(this.accumulator / (1000/this.updateRate))
    }

    insertDrawable(drawable, layer){
        this.drawablesLayers.insertDrawable(drawable, layer)

        if (drawable.rect.isKinematic()){
            this.kinematicRects.push(drawable.rect)
        }
    }

    interpolateKinematics(alpha){
        this.kinematicRects.forEach(rect => {
            rect.interpolate(alpha)
        })
    }

    getCameraPos(){
        return [this.camera.left, this.camera.top]
    }
}