function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class Rect {
    left
    top
    width
    height

    constructor(){
        this.setPosition(0, 0)
        this.setSize(0, 0)
    }

    setPosition(left, top){
        this.left = left
        this.top = top
    }

    setSize(width, height){
        this.width = width
        this.height = height
    }

    containsPoint(x, y){
        if(
        this.left <= x && x <= this.left + this.width &&
        this.top <= y && y <= this.top + this.height
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
    prev
    curr

    constructor(){
        super()
        this.prev = new Rect()
        this.prev.setPosition(this.left, this.top)
        this.prev.setSize(this.width, this.height)

        this.curr = new Rect()
        this.curr.setPosition(this.left, this.top)
        this.curr.setSize(this.width, this.height)
    }

    setCurrentPos(left, top, teleport){
        this.curr.setPosition(left, top)

        if (teleport) {
            this.setPosition(left, top)
            this.prev.setPosition(left, top)
        }
    }

    updatePrevPos(){
        this.prev.setPosition(this.curr.left, this.curr.top)
    }

    interpolate(alpha){
        let posX = lerp(this.prev.left, this.curr.left, alpha) 
        let posY = lerp(this.prev.top, this.curr.top, alpha)
        this.setPosition(posX, posY)
    }

    isKinematic(){
        return true
    }
} 