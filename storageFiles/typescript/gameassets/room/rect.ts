function lerp(a: number, b: number, alpha: number): number{
    return a + alpha * (b - a);
}

interface AbstractRect {
    position: {
        x: number,
        y: number,
    }
    size: {
        x: number,
        y: number,
    }
}

function isAbstractRect(obj: any): obj is AbstractRect {
    return "position" in obj && "size" in obj
}

class Rect {
    protected position: Vector2;
    protected size: Vector2;

    constructor(abstractRect?: AbstractRect){
        this.position = new Vector2(0, 0);
        this.size = new Vector2(0, 0);

        if(abstractRect){
            this.setPosition(abstractRect.position.x, abstractRect.position.y);
            this.setSize(abstractRect.size.x, abstractRect.size.y);
        }
    }

    setPosition(x: number, y: number){
        this.position.setPosition(x, y);
    }

    setSize(x: number, y: number){
        this.size.setPosition(x, y);
    }

    extend(extX: number, extY: number){
        this.size.x += Math.abs(extX)
        this.size.y += Math.abs(extY)
        if (extX < 0){
            this.position.x -= Math.abs(extX)
        }
        if (extY < 0){
            this.position.y -= Math.abs(extY)
        }
    }

    containsPoint(x: number, y: number){
        let pos = this.position;
        let size = this.size;

        if(
        x >= pos.x && x <= pos.x + size.x &&
        y >= pos.y && y <= pos.y + size.y
        ){
            return true;
        }
        return false;
    }

    intersects(rect: Rect){
        if (this.position.x+this.size.x <= rect.position.x ||
        this.position.x >= rect.position.x+rect.size.x ||
        this.position.y+this.size.y <= rect.position.y ||
        this.position.y >= rect.position.y+rect.size.y){
            return false
        }

        return true
    }

    getPosition(){
        return new Vector2(this.position.x, this.position.y);
    }

    getSize(){
        return new Vector2(this.size.x, this.size.y);
    }
}