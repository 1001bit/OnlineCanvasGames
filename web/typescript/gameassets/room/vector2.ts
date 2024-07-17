function lerpVector2(v1: Vector2, v2: Vector2, a: number){
    return new Vector2(
        v1.x + a * (v2.x - v1.x), 
        v1.y + a * (v2.y - v1.y)
    )
}

class Vector2 {
    x: number;
    y: number;

    constructor(x: number, y: number){
        this.x = x;
        this.y = y;
    }

    setPosition(x: number, y: number){
        this.x = x;
        this.y = y;
    }
}