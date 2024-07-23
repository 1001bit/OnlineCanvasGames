interface AbstractPlayer extends AbstractRect {

}

function isAbstractPlayer(obj: any) : obj is AbstractPlayer{
    return isAbstractRect(obj)
}

class InterpolatedPlayer extends Rect {
    startPosition: Vector2;
    targetPosition: Vector2;

    constructor(abstract: AbstractPlayer){
        super(abstract)

        this.startPosition = this.getPosition()
        this.targetPosition = this.getPosition()
    }

    setTargetPos(x: number, y: number, teleport?: boolean){
        this.targetPosition.setPosition(x, y)

        if(teleport){
            this.setPosition(x, y)
            this.startPosition.setPosition(x, y)
        }
    }

    updateStartPos(){
        this.startPosition.setPosition(this.targetPosition.x, this.targetPosition.y)
    }

    interpolate(alpha: number){
        const newPos = lerpVector2(this.startPosition, this.targetPosition, alpha)
        this.setPosition(newPos.x, newPos.y)
    }
}

class KinematicPlayer extends InterpolatedPlayer {
    velocity: Vector2;

    collisionVertical: Direction
    collisionHorizontal: Direction

    constructor(abstract: AbstractPlayer){
        super(abstract)

        this.velocity = new Vector2(0, 0)
        this.collisionHorizontal = Direction.None
        this.collisionVertical = Direction.None
    }

    control(speed: number, jump: number, controls: Controls){
        if(controls.isHeld("left")){
            this.velocity.x -= speed
        }

        if(controls.isHeld("right")){
            this.velocity.x += speed
        }

        if(controls.isHeld("jump") && this.collisionVertical == Direction.Down){
            this.velocity.y -= jump
        }
    }

    applyGravity(force: number, dt: number){
        this.velocity.y += force * dt
    }

    applyFriction(force: number){
        this.velocity.x *= force
	    // this.velocity.y *= force
    }

    setCollisionDir(dir: Direction){
        if (dir == Direction.Down || dir == Direction.Up){
            this.collisionVertical = dir
        } else if (dir == Direction.Left || dir == Direction.Right){
            this.collisionHorizontal = dir
        } else {
            this.collisionHorizontal = dir
            this.collisionVertical = dir
        }
    }
}