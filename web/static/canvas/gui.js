class Gui{
    constructor(layers){
        this.drawablesLayers = new DrawableLayers(layers)

        this.tickRate = 60
        this.tickInterval = setInterval(() => this.tick(), 1000/this.tickRate)
        this.timer = new DeltaTimer()
    }

    tick(){

    }

    draw(ctx){
        this.drawablesLayers.draw(ctx)
    }

    insertDrawable(drawable, layer){
        this.drawablesLayers.insertDrawable(drawable, layer)
    }
}