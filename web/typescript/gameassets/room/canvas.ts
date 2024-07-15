class GameCanvas {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;
    
    layers: Array<Map<number, Drawable>>;
    drawables: Map<number, Drawable>;

    mousePos: Vector2;

    backgroundColor: string;

    constructor(canvasID: string, layersCount: number) {
        this.canvas = <HTMLCanvasElement> document.getElementById(canvasID);

        const ctx = this.canvas.getContext("2d");
        if(!ctx){
            throw new Error("Failed to get context");
        }
        this.ctx = ctx;
        
        this.layers = new Array();
        for (let i = 0; i < layersCount; i++){
            this.layers.push(new Map());
        }
        this.drawables = new Map();

        this.mousePos = new Vector2(0, 0);

        this.backgroundColor = RGB(0, 0, 0);

        this.resize();
        window.addEventListener('resize', () => this.resize(), false);

        this.canvas.addEventListener("mousemove", e => {
            this.updateMousePos(e);
        })
    }

    stop(){
        this.canvas.remove();
    }

    resize(){
        const canvas = this.canvas;

        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top;

        this.draw();
    }

    insertDrawable(drawable: Drawable, layerNum: number, id: number){
        if(this.drawables.has(id)){
            return
        }

        this.drawables.set(id, drawable);

        const layer = this.layers[layerNum]
        if(layer){
            layer.set(id, drawable);
        }
    }

    deleteDrawable(id: number){
        this.drawables.delete(id);
        this.layers.forEach(layer => {
            layer.delete(id);
        })
    }

    draw(){
        const ctx = this.ctx;

        this.clear();

        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                drawable.draw(ctx);
            });
        });
    }

    clear(){
        const ctx = this.ctx;
        const canvas = this.canvas;

        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
    }

    setBackgroundColor(color: string){
        this.backgroundColor = color;
    }

    updateMousePos(e: MouseEvent){
        let rect = this.canvas.getBoundingClientRect();
        let x = e.clientX - rect.left;
        let y = e.clientY - rect.top;
        this.mousePos.setPosition(x, y);
    }

    getMousePos(){
        return this.mousePos;
    }

    getDrawable(id: number){
        return this.drawables.get(id);
    }
}