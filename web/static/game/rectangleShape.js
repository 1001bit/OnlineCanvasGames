class RectangleShape {
    rect
    color

    constructor(width, height, kinematic){
        this.rect = kinematic ? new KinematicRect() : new Rect()
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