class Level{
    constructor(layers){
        this.drawablesLayers = new DrawableLayers(layers)
        this.camera = new KinematicRect()
        this.kinematicRects = [this.camera]

        this.accumulator = 0
        this.tickRate = 60
        this.tickInterval = setInterval(() => this.tick(), 1000/this.tickRate)
        this.timer = new DeltaTimer()
    }

    tick(){
        let dt = this.timer.getDeltaTime()

        this.update(dt)
        this.draw()
    }

    update(){

    }

    draw(){
        
    }
}