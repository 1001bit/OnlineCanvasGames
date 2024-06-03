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
}

class KinematicRect extends Rect {
    constructor(){
        super()
        this.previousPos = new Vector2(0, 0)
        this.targetPos = new Vector2(0, 0)
    }

    setTargetPos(x, y, teleport){
        this.targetPos.setPosition(x, y)

        if (teleport) {
            this.setPosition(x, y)
            this.updatePrevPos()
        }
    }

    updatePrevPos(){
        this.previousPos.setPosition(this.targetPos.x, this.targetPos.y)
    }

    interpolate(alpha){
        let prev = this.previousPos
        let targ = this.targetPos

        let posX = lerp(prev.x, targ.x, alpha) 
        let posY = lerp(prev.y, targ.y, alpha)
        this.setPosition(posX, posY)
    }

    isKinematic(){
        return true
    }
} 