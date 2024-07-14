class DrawableText extends Drawable {
    string: string;
    color: string;
    font: string;
    size: number;
    
    position: Vector2;

    constructor(string: string, size: number){
        super();

        this.string = string;
        this.color = RGB(255, 255, 255);
        this.font = "serif";
        this.size = size;
        
        this.position = new Vector2(0, 0);
    }

    setPosition(x: number, y: number){
        this.position.setPosition(x, y);
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

    setSize(size: number){
        this.size = size;
    }

    override draw(ctx: CanvasRenderingContext2D){
        ctx.fillStyle = this.color;
        ctx.font = `${this.size}px ${this.font}`;
        // adding size to y because text's origin is located on the bottom
        ctx.fillText(this.string, this.position.x, this.position.y + this.size);
    }
}