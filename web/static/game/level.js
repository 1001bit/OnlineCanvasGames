class Level {
    constructor(drawables){
        this.drawables = drawables
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
}