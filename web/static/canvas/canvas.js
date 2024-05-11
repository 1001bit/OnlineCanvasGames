function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class GameCanvas {
    canvas
    ctx

    drawables
    kinematicRects

    backgroundColor

    updateRate
    updateInterval
    accumulator

    drawRate
    drawInterval
    lastDraw

    mousePos

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        this.drawables = []
        this.accumulator = 0
        this.lastDraw = Date.now()
        this.kinematicRects = []

        this.updateRate = 30
        this.drawRate = 60

        this.mousePos = [0, 0]

        this.setBackgroundColor(RGB(0, 0, 0))

        window.addEventListener('resize', () => this.resize(), false);
        $(document).mousemove(e => {
            this.updateMousePos(e)
        })
    }

    start(){
        $("header").hide()
        this.setCanvasVisibility(true)
        this.resize()

        clearInterval(this.drawInterval)
        this.drawInterval = setInterval(() => this.draw(), 1000/this.drawRate)

        clearInterval(this.updateInterval)
        this.updateInterval = setInterval(() => this.update(), 1000/this.updateRate)
    }

    stop(){
        $("header").show()
        this.setCanvasVisibility(false)

        clearInterval(this.interval)
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

    insertNewDrawable(drawable, hasKinematicRect){
        this.drawables.push(drawable)

        if (hasKinematicRect){
            this.kinematicRects.push(drawable.rect)
        }
    }

    interpolateKinematics(){
        let now = Date.now()
        let dt = now - this.lastDraw
        this.lastDraw = now

        this.accumulator += dt
        while(this.accumulator >= 1000/this.updateRate){
            this.updateKinematics()
            this.accumulator -= 1000/this.updateRate
        }

        let alpha = this.accumulator / (1000/this.updateRate)

        this.kinematicRects.forEach(rect => {
            let posX = lerp(rect.prev.left, rect.curr.left, alpha) 
            let posY = lerp(rect.prev.top, rect.curr.top, alpha)
            rect.setPosition(posX, posY)
        })
    }

    updateKinematics(){
        this.kinematicRects.forEach(rect => {
            rect.updatePrevPos()
        })
    }

    draw(){
        this.interpolateKinematics()

        this.clear()
        this.drawables.forEach(drawable => {
            drawable.draw(this.ctx)
        })
    }

    clear(){
        const ctx = this.ctx
        const canvas = this.canvas

        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor
        ctx.fillRect(0, 0, canvas.width, canvas.height)
    }

    update = () => {}

    setBackgroundColor(color){
        this.backgroundColor = color
    }

    updateMousePos(e){
        let rect = this.canvas.getBoundingClientRect()
        let x = e.clientX - rect.left
        let y = e.clientY - rect.top
        this.mousePos = [x, y]
    }

    getMousePos(){
        return this.mousePos
    }
}