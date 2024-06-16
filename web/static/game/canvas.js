function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class GameCanvas {
    constructor(canvasID, layersCount) {
        this.active = true

        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        
        this.layers = []
        for (let i = 0; i < layersCount; i++){
            this.layers.push(new Map())
        }

        this.camera = new KinematicRect()

        this.mousePos = new Vector2(0, 0)

        this.resize()

        window.addEventListener('resize', () => this.resize(), false);

        $(document).mousemove(e => {
            this.updateMousePos(e)
        })
    }

    stop(){
        this.active = false
        this.clear()
        this.canvas.remove()
    }

    resize (){
        if(!this.active){
            return
        }

        const canvas = this.canvas

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top

        this.draw()
    }

    insertDrawable(drawable, layer, id){
        this.layers[layer].set(id, drawable)
    }

    deleteDrawable(id){
        this.layers.forEach(layer => {
            layer.delete(id)
        })
    }

    setCameraPos(x, y){
        this.camera.setTargetPos(x, y)
    }

    draw(){
        const ctx = this.ctx
        this.clear()

        const cameraPos = this.camera.getPosition()

        ctx.save()
        ctx.translate(-cameraPos.x, -cameraPos.y) // for some reason, it has to be a negative value

        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                console.log(1)
                drawable.draw(ctx)
            })
        })

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

    updateMousePos(e){
        let rect = this.canvas.getBoundingClientRect()
        let x = e.clientX - rect.left
        let y = e.clientY - rect.top
        this.mousePos.setPosition(x, y)
    }

    getMousePos(){
        return this.mousePos
    }

    getLevelMousePos(){
        let cameraPos = this.camera.getPosition()
        let mousePos = this.mousePos

        return new Vector2(cameraPos.x + mousePos.x, cameraPos.y + mousePos.y)
    }

    getCameraPos(){
        return this.camera.getPosition()
    }
}