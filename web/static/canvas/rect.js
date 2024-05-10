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
}

class RectangleShape {
    rect
    color

    constructor(width, height){
        this.rect = new Rect()
        this.rect.setSize(width, height)
        this.rect.setPosition(0, 0)

        this.setColor(RGB(255, 255, 255))
    }

    setPosition(left, top){
        this.rect.setPosition(left, top)
    }

    setSize(width, height){
        this.rect.setSize(width, height)
    }

    setColor(fillColor){
        this.color = fillColor
    }

    draw(ctx){
        ctx.fillStyle = this.color
        ctx.fillRect(this.rect.left, this.rect.top, this.rect.width, this.rect.height)
    }
}