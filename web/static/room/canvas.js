function RGB(r, g, b){
    return `rgb(${r} ${g} ${b})`
}

class Rect {
    left
    top
    width
    height

    constructor(){
        this.left = 0
        this.top = 0
        this.width = 0
        this.height = 0
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
    fillColor

    constructor(width, height){
        this.rect = new Rect()
        this.rect.setSize(width, height)
        this.rect.setPosition(0, 0)

        this.fillColor = "rgb(0 0 0)"
    }

    setPosition(left, top){
        this.rect.setPosition(left, top)
    }

    setSize(width, height){
        this.rect.setSize(width, height)
    }

    setColor(fillColor){
        this.fillColor = fillColor
    }

    draw(ctx){
        ctx.fillStyle = this.fillColor
        ctx.fillRect(this.rect.left, this.rect.top, this.rect.width, this.rect.height)
    }
}

class GameCanvas {
    canvas
    ctx
    drawables
    interval

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        this.drawables = []

        window.addEventListener('resize', () => this.resize(), false);
    }

    start(){
        $("header").hide()
        this.canvas.style.display = "block" 
        this.resize()

        const rate = 60
        this.interval = setInterval(() => this.draw(), 1000/rate)
    }

    stop(){
        $("header").show()
        this.canvas.style.display = "none" 

        clearInterval(this.interval)
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
        ctx.fillStyle = "rgb(0 0 0)"
        ctx.fillRect(0, 0, canvas.width, canvas.height)
    }
}