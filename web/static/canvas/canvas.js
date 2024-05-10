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

    getMousePosition(e){
        let rect = this.canvas.getBoundingClientRect()
        let x = e.clientX - rect.left
        let y = e.clientY - rect.top
        return {x, y}
    }
}