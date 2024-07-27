class GameCanvas {
    private canvas: HTMLCanvasElement;
    private ctx: CanvasRenderingContext2D;
    
    private layers: Array<Map<number, Drawable>>;
    private drawables: Map<number, Drawable>;

    private mousePos: Vector2;

    private backgroundColor: string;

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

        this.canvas.addEventListener("click", e => {
            this.onMouseClick(e);
        })
    }

    private resize(){
        const canvas = this.canvas;

        canvas.width = window.innerWidth;
        canvas.height = window.innerHeight - canvas.getBoundingClientRect().top;

        this.draw();
    }

    private clear(){
        const ctx = this.ctx;
        const canvas = this.canvas;

        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.fillStyle = this.backgroundColor;
        ctx.fillRect(0, 0, canvas.width, canvas.height);
    }

    private updateMousePos(e: MouseEvent){
        let rect = this.canvas.getBoundingClientRect();
        let x = e.clientX - rect.left;
        let y = e.clientY - rect.top;
        this.mousePos.setPosition(x, y);
    }

    onMouseClick = (_e: MouseEvent) => {}

    stop(){
        this.canvas.remove();
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

    draw(viewCenter?: Vector2){
        const ctx = this.ctx;
        ctx.save();

        this.clear();

        if(viewCenter){
            ctx.translate(this.canvas.width/2 - viewCenter.x, this.canvas.height/2 - viewCenter.y)
        }

        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                drawable.draw(ctx);
            });
        });

        ctx.restore();
    }

    setBackgroundColor(color: string){
        this.backgroundColor = color;
    }

    getMousePos(){
        return this.mousePos;
    }

    getDrawable(id: number){
        return this.drawables.get(id);
    }
}