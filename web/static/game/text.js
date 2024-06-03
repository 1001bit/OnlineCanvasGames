class Text {
    constructor(string, size){
        this.string = string
        this.color = RGB(255, 255, 255)
        this.font = "serif"
        this.size = size
        
        this.position = new Vector2()
    }

    setPosition(x, y){
        this.position.setPosition(x, y)
    }

    setString(string){
        this.string = string
    }

    setColor(color){
        this.color = color
    }

    setFont(font){
        this.font = font
    }

    setSize(size){
        this.size = size
    }

    draw(ctx){
        ctx.fillStyle = this.color
        ctx.font = `${this.size}px ${this.font}`
        // adding size to y because text's origin is located on the bottom
        ctx.fillText(this.string, this.position.x, this.position.y + this.size) 
    }
}