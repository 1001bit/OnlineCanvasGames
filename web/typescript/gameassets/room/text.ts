class DrawableText extends Drawable {
    string: string;
    color: string;
    font: string;
    fontSize: number;

    position: Vector2;

    constructor(string: string, fontSize: number){
        super();

        this.string = string;
        this.color = RGB(255, 255, 255);
        this.fontSize = fontSize
        this.font = "serif";
        this.position = new Vector2(0, 0)
    }

    setString(string: string){
        this.string = string;
    }

    setColor(color: string){
        this.color = color;
    }

    setFont(font: string){
        this.font = font;
    }

    setFontSize(fontSize: number){
        this.fontSize = fontSize
    }

    setPosition(x: number, y: number){
        this.position.setPosition(x, y)
    }

    override draw(ctx: CanvasRenderingContext2D){
        ctx.fillStyle = this.color;
        ctx.font = `${this.fontSize}px ${this.font}`;

        // adding height to y because text's origin is located on the bottom
        const metrics = ctx.measureText(this.string)
        const height = metrics.actualBoundingBoxAscent + metrics.actualBoundingBoxDescent
    
        ctx.fillText(this.string, this.position.x, this.position.y + height);
    }
}