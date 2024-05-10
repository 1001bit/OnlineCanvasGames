class Text {
    string
    color
    font
    size

    x
    y

    constructor(string, size){
        this.setString(string)
        this.setColor(RGB(255, 255, 255))
        this.setFont("serif")
        this.setSize(size)
        this.setPosition(0, 0)
    }

    setPosition(x, y){
        this.x = x
        this.y = y
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
        ctx.fillText(this.string, this.x, this.y + this.size) 
    }
}