class RectangleShape extends Drawable {
    color: string;

    constructor(width: number, height: number){
        super();

        this.setSize(width, height)

        this.color = RGB(255, 255, 255);
    }

    setColor(color: string){
        this.color = color;
    }

    override draw(ctx: CanvasRenderingContext2D){
        let pos = this.position;
        let size = this.size;

        ctx.fillStyle = this.color;
        ctx.fillRect(pos.x, pos.y, size.x, size.y);
    }
}