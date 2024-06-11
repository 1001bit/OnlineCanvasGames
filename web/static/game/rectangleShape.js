class RectangleShape {
    constructor(width, height, kinematic){
        this.rect = kinematic ? new KinematicRect() : new Rect()
        this.rect.setSize(width, height)
        this.rect.setPosition(0, 0)

        this.color = RGB(255, 255, 255) 
    }

    setSize(width, height){
        this.rect.setSize(width, height)
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