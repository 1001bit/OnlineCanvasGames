function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class Rect {
    constructor(){
        this.position = new Vector2(0, 0)
        this.size = new Vector2(0, 0)
    }

    setPosition(x, y){
        this.position.setPosition(x, y)
    }

    setSize(x, y){
        this.size.setPosition(x, y)
    }

    containsPoint(x, y){
        let pos = this.position
        let size = this.size

        if(
        x >= pos.x && x <= pos.x + size.x &&
        y >= pos.y && y <= pos.y + size.y
        ){
            return true
        }
        return false
    }

    isKinematic(){
        return false
    }

    getPosition(){
        return this.position
    }

    getSize(){
        return this.size
    }
}

class KinematicRect extends Rect {
    constructor(){
        super()
        this.velocity = new Vector2(0, 0)
        this.serverPos = new Vector2(0, 0)
    }

    setVelocity(x, y){
        this.velocity.x = x
        this.velocity.y = y
    }

    setServerPos(x, y){
        this.serverPos.x = x
        this.serverPos.y = y
    }

    applyVelToPos(dt){
        let posX = this.position.x + this.velocity.x * dt
        let posY = this.position.y + this.velocity.y * dt
        this.setPosition(posX, posY)
    }

    pullToServerPos(dt){
        const approachFactor = 0.8

        let blend = Math.pow(approachFactor, dt)
        this.position.x = lerp(this.position.x, this.serverPos.x, blend)
        this.position.y = lerp(this.position.y, this.serverPos.y, blend)
    }

    isKinematic(){
        return true
    }
} 