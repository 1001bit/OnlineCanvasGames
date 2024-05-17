function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class GameCanvas {
    gameUpdate = () => {}

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        
        this.gui = new Gui(2)
        this.level = new Level(2)

        this.setBackgroundColor(RGB(0, 0, 0))

        this.mousePos = [0, 0]

        this.tickRate = 60
        this.tickInterval = setInterval(() => this.tick(), 1000/this.tickRate)

        this.setCanvasVisibility(true)

        window.addEventListener('resize', () => this.resize(), false);

        $(document).mousemove(e => {
            this.updateMousePos(e)
        })
    }

    stop(){
        this.setCanvasVisibility(false)
        clearInterval(this.tickInterval)
    }

    setCanvasVisibility(visibility){
        if (visibility){
            this.canvas.style.display = "block"
            $("header").hide()
            this.resize()
            return
        }
        this.canvas.style.display = "none"
        $("header").show()
    }

    resize(){
        const canvas = this.canvas

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top
        this.draw()
    }

    setCameraPos(x, y){
        this.camera.setCurrentPos(x, y)
    }

    tick(){
        this.draw()
    }

    draw(){
        const ctx = this.ctx

        this.clear()

        this.level.draw(ctx)
        this.gui.draw(ctx)
    }

    clear(){
        const ctx = this.ctx
        const canvas = this.canvas

        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor
        ctx.fillRect(0, 0, canvas.width, canvas.height)
    }

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

    getLevelMousePos(){
        let [mx, my] = this.mousePos
        let [vx, vy] = this.level.getCameraPos()

        return [vx + mx, vy + my]
    }
}