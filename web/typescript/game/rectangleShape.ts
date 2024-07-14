class RectangleShape extends Drawable {
    rect: Rect;
    color: string;

    constructor(width: number, height: number){
        super();

        this.rect = new Rect();
        this.rect.setSize(width, height);
        this.rect.setPosition(0, 0);

        this.color = RGB(255, 255, 255);
    }

    setSize(width: number, height: number){
        this.rect.setSize(width, height);
    }

    setPosition(x: number, y: number){
        this.rect.setPosition(x, y);
    }

    setColor(color: string){
        this.color = color;
    }

    override draw(ctx: CanvasRenderingContext2D){
        let pos = this.rect.position;
        let size = this.rect.size;

        ctx.fillStyle = this.color;
        ctx.fillRect(pos.x, pos.y, size.x, size.y);
    }
}