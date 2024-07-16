class RectangleShape extends Drawable {
    color: string;
    rect: Rect;

    constructor(rect?: Rect){
        super();

        if(rect){
            this.rect = rect
        } else {
            this.rect = new Rect()
        }

        this.color = RGB(255, 255, 255);
    }

    setSize(x: number, y: number){
        this.rect.setSize(x, y)
    }

    setPosition(x: number, y: number){
        this.rect.setPosition(x, y)
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