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

    detectHorizontalCollision(block: Block, dtMs: number){
        if(this.velocity.x == 0){
            return Direction.None
        }

        const playerPath = new Rect(this)
        playerPath.setPosition(this.targetPosition.x, this.targetPosition.y)
        playerPath.extend(this.velocity.x * dtMs, 0)

        if(!playerPath.intersects(block)){
            return Direction.None
        }

        if(this.velocity.x > 0){
            return Direction.Right
        } else {
            return Direction.Left
        }
    }

    detectVerticalCollision(block: Block, dtMs: number){
        if(this.velocity.y == 0){
            return Direction.None
        }

        const playerPath = new Rect(this)
        playerPath.setPosition(this.targetPosition.x, this.targetPosition.y)
        playerPath.extend(0, this.velocity.y * dtMs)

        if(!playerPath.intersects(block)){
            return Direction.None
        }

        if(this.velocity.y > 0){
            return Direction.Down
        } else {
            return Direction.Up
        }
    }

    resolveCollision(block: Block, dir: Direction){
        if(dir == Direction.None){
            return
        }

        this.setCollisionDir(dir)

        switch(dir){
            case Direction.Up:
                this.velocity.y = 0
                this.targetPosition.y = block.getPosition().y + block.getSize().y
                break
            case Direction.Down:
                this.velocity.y = 0
                this.targetPosition.y = block.getPosition().y - this.getSize().y
                break

            case Direction.Left:
                this.velocity.x = 0
                this.targetPosition.x = block.getPosition().x + block.getSize().x
                break
            case Direction.Right:
                this.velocity.x = 0
                this.targetPosition.x = block.getPosition().x - this.getSize().x
                break
        }
    }
}