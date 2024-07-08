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

    getPosition(){
        return this.position
    }

    getSize(){
        return this.size
    }
}