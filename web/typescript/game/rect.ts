function lerp(a: number, b: number, alpha: number): number{
    return a + alpha * (b - a);
}

class Rect {
    position: Vector2;
    size: Vector2;

    constructor(){
        this.position = new Vector2(0, 0);
        this.size = new Vector2(0, 0);
    }

    setPosition(x: number, y: number){
        this.position.setPosition(x, y);
    }

    setSize(x: number, y: number){
        this.size.setPosition(x, y);
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

    getPosition(){
        return this.position;
    }

    getSize(){
        return this.size;
    }
}