class RectangleShape {
    constructor(width, height){
        this.rect = new Rect()
        this.rect.setSize(width, height)
        this.rect.setPosition(0, 0)

        this.color = RGB(255, 255, 255) 
    }

    setSize(width, height){
        this.rect.setSize(width, height)
    }

    setPosition(left, top){
        this.rect.setPosition(left, top)
    }

    setColor(color){
        this.color = color
    }

    draw(ctx){
        let pos = this.rect.position
        let size = this.rect.size

        ctx.fillStyle = this.color
        ctx.fillRect(pos.x, pos.y, size.x, size.y)
    }
}