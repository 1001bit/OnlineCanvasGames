class Text {
    string
    color
    font
    size

    rect

    constructor(string, size){
        this.setString(string)
        this.setColor(RGB(255, 255, 255))
        this.setFont("serif")
        this.setSize(size)
        
        this.rect = new Rect()
    }

    setPosition(x, y){
        this.rect.setPosition(x, y)
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
        ctx.fillText(this.string, this.rect.left, this.rect.top + this.size) 
    }
}