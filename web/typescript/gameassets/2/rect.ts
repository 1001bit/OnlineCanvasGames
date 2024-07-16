function isRect(obj: any) : obj is Rect{
    return "position" in obj && "size" in obj
}

function isKinematicRect(obj: any) : obj is KinematicRect{
    return isRect(obj) && "velocity" in obj
}

class KinematicRect extends Rect {
    velocity: Vector2;

    constructor(rect?: Rect){
        super(rect);

        this.velocity = new Vector2(0, 0);
    }

    setVelocity(x: number, y: number){
        this.velocity.setPosition(x, y)
    }
}