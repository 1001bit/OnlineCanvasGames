function RGB(r, g, b){
    return `rgb(${r} ${g} ${b})`
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

class GameCanvas {
    canvas
    ctx
    drawables
    backgroundColor
    interval

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        this.drawables = []

        this.setBackgroundColor(RGB(0, 0, 0))

        window.addEventListener('resize', () => this.resize(), false);
    }

    start(){
        $("header").hide()
        this.setCanvasVisibility(true)
        this.resize()
        this.setDrawRate(60)
    }

    stop(){
        $("header").show()
        this.setCanvasVisibility(false)

        clearInterval(this.interval)
    }

    setDrawRate(rate){
        clearInterval(this.interval)
        this.interval = setInterval(() => this.draw(), 1000/rate)
    }

    setCanvasVisibility(visibility){
        this.canvas.style.display = visibility ? "block" : "none"
    }

    resize(){
        const canvas = this.canvas

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top
        this.draw()
    }

    insertNewDrawable(drawable){
        this.drawables.push(drawable)
    }

    draw(){
        const ctx = this.ctx
        const canvas = this.canvas

        this.cleanCanvas()

        this.drawables.forEach(drawable => {
            drawable.draw(ctx)
        })
    }

    cleanCanvas(){
        const ctx = this.ctx
        const canvas = this.canvas

        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor
        ctx.fillRect(0, 0, canvas.width, canvas.height)
    }

    setBackgroundColor(color){
        this.backgroundColor = color
    }
}