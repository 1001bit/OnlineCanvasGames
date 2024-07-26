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

    interpolateBetween(v1: Vector2, v2: Vector2, a: number){
        this.x = lerp(v1.x, v2.x, a)
        this.y = lerp(v1.y, v2.y, a)
    }
}