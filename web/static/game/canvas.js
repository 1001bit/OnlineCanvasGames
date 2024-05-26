function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class GameCanvas {
    gameUpdate = () => {}

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        
        this.drawablesLayers = new DrawablesLayers(1)
        this.camera = new KinematicRect()
        this.kinematicRects = [this.camera]

        this.mousePos = [0, 0]

        this.drawRate = 60
        this.drawInterval = setInterval(() => this.draw(), 1000/this.drawRate)

        this.resize()

        window.addEventListener('resize', () => this.resize(), false);

        $(document).mousemove(e => {
            this.updateMousePos(e)
        })
    }

    stop(){
        clearInterval(this.drawInterval)
        this.setBackgroundColor(RGB(0, 0, 0))
        this.clear()
    }

    resize = () => {
        const canvas = this.canvas

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top
        this.draw()
    }

    setCameraPos(x, y){
        this.camera.setCurrentPos(x, y)
    }

    draw(){
        const ctx = this.ctx
        this.clear()

        ctx.save()
        ctx.translate(-this.camera.left, -this.camera.top) // for some reason, it has to be a negative value

        // there is no need for interpolation
        this.kinematicRects.forEach(rect => {
            rect.updatePrevPos()
        })
        this.interpolateKinematics(1)

        this.drawablesLayers.draw(ctx)

        ctx.restore()
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

    setLayersCount(layers){
        this.drawablesLayers = new DrawablesLayers(layers)
    }

    insertDrawable(drawable, layer){
        this.drawablesLayers.insertDrawable(drawable, layer)
        if (drawable.rect.isKinematic()){
            this.kinematicRects.push(drawable.rect)
        }
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
        let [vx, vy] = [this.camera.left, this.camera.top]

        return [vx + mx, vy + my]
    }

    interpolateKinematics(alpha){
        this.kinematicRects.forEach(rect => {
            rect.interpolate(alpha)
        })
    }
}