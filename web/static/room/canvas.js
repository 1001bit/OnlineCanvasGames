class GameCanvas {
    canvas
    ctx

    constructor(canvasID) {
        this.canvas = document.getElementById(canvasID)
        this.ctx = this.canvas.getContext("2d")
        window.addEventListener('resize', () => this.resize(), false);
    }

    resize(){
        const canvas = this.canvas

        canvas.width = window.innerWidth
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top - 6
        this.draw()
    }

    draw(){
        const ctx = this.ctx

        ctx.fillStyle = "rgb(255 0 0)"
        ctx.fillRect(10, 10, 50, 50)
    }
}