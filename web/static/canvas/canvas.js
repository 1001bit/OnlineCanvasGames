function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class GameCanvas {
    canvas
    ctx

    levelDrawables
    guiDrawables
    camera

    kinematicRects

    backgroundColor

    updateRate
    accumulator

    tickRate
    tickInterval
    lastTick

    mousePos

    gameUpdate = () => {}

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        
        this.levelDrawableLayers = new DrawableLayers(2)
        this.guiDrawableLayers = new DrawableLayers(2)
        this.camera = new KinematicRect()

        this.kinematicRects = [this.camera]

        this.setBackgroundColor(RGB(0, 0, 0))

        this.updateRate = 20
        this.accumulator = 0

        this.tickRate = 60
        this.lastTick = Date.now()

        this.mousePos = [0, 0]

        window.addEventListener('resize', () => this.resize(), false);

        $(document).mousemove(e => {
            this.updateMousePos(e)
        })
    }

    start(){
        $("header").hide()
        this.setCanvasVisibility(true)
        this.resize()

        clearInterval(this.tickInterval)
        this.tickInterval = setInterval(() => this.tick(), 1000/this.tickRate)
    }

    stop(){
        $("header").show()
        this.setCanvasVisibility(false)

        clearInterval(this.tickInterval)
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

    setCameraPos(x, y){
        this.camera.setCurrentPos(x, y)
    }

    insertLevelDrawable(drawable, layer){
        this.levelDrawableLayers.insertDrawable(drawable, layer)

        if (drawable.rect.isKinematic()){
            this.kinematicRects.push(drawable.rect)
        }
    }

    insertGuiDrawable(drawable, layer){
        this.guiDrawableLayers.insertDrawable(drawable, layer)

        if (drawable.rect.isKinematic()){
            this.kinematicRects.push(drawable.rect)
        }
    }

    tick(){
        let now = Date.now()
        let dt = now - this.lastTick
        this.lastTick = now

        this.update(dt)
        this.draw()
    }

    update(dt){
        this.accumulator += dt

        while(this.accumulator >= 1000/this.updateRate){
            this.accumulator -= 1000/this.updateRate

            this.kinematicRects.forEach(rect => {
                rect.updatePrevPos()
            })
            this.gameUpdate()
        }

        this.interpolateKinematics(this.accumulator / (1000/this.updateRate))
    }

    interpolateKinematics(alpha){
        this.kinematicRects.forEach(rect => {
            rect.interpolate(alpha)
        })
    }

    draw(){
        const ctx = this.ctx

        this.clear()

        ctx.save()
        ctx.translate(-this.camera.left, -this.camera.top) // for some reason, it has to be a negatile value

        this.levelDrawableLayers.draw(ctx)

        ctx.restore()

        this.guiDrawableLayers.draw(ctx)
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

    getGuiMousePos(){
        return this.mousePos
    }

    getLevelMousePos(){
        let [mx, my] = this.mousePos
        let [vx, vy] = [this.camera.left, this.camera.top]

        return [vx + mx, vy + my]
    }
}