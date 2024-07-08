function lerp(a, b, alpha){
    return a + alpha * (b - a)
}

class GameCanvas {
    constructor(canvasID, layersCount) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        
        this.layers = new Array()
        for (let i = 0; i < layersCount; i++){
            this.layers.push(new Map())
        }
        this.drawables = new Map()

        this.mousePos = new Vector2(0, 0)

        this.resize()
        window.addEventListener('resize', () => this.resize(), false);

        $(document).mousemove(e => {
            this.updateMousePos(e)
        })
    }

    stop(){
        this.canvas.remove()
    }

    resize(){
        const canvas = this.canvas

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top

        this.draw()
    }

    insertDrawable(drawable, layer, id){
        this.drawables.set(id, drawable)
        this.layers[layer].set(id, drawable)
    }

    deleteDrawable(id){
        this.drawables.delete(id)
        this.layers.forEach(layer => {
            layer.delete(id)
        })
    }

    draw(){
        const ctx = this.ctx

        this.clear()

        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                drawable.draw(ctx)
            })
        })
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

    drawableExists(id){
        return this.drawables.has(id)
    }

    getDrawable(id){
        return this.drawables.get(id)
    }
}