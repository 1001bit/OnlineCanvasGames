class Level {
    constructor(canvas, drawables){
        this.drawables = drawables
        this.canvas = canvas

        this.serverTPS = 0
        this.clientTPS = 0

        this.timer = new DeltaTimer()
    }

    handleLevelMessage = (body) => {
        for (const key in body){
            if(!(key in this.drawables)){
                continue
            }
            
            this.drawables[key].setPosition(body[key].x, body[key].y)
            this.drawables[key].setSize(body[key].w, body[key].h)
        }
    }

    setTPS(client, server){
        this.clientTPS = client
        this.serverTPS = server

        this.timer.getDeltaTime()
        setInterval(() => {this.tick()}, 1000/this.clientTPS)
    }

    update(dt){

    }

    tick(){
        let dt = this.timer.getDeltaTime()
        this.update(dt)

        this.canvas.draw()
    }
}