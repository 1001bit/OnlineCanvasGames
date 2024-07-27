interface AbstractPlayer extends AbstractRect {

}

function isAbstractPlayer(obj: any) : obj is AbstractPlayer{
    return isAbstractRect(obj)
}

class InterpolatedPlayer extends Rect {
    protected startPosition: Vector2;
    protected targetPosition: Vector2;

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

    moveTargetPos(x: number, y: number){
        this.targetPosition.x += x
        this.targetPosition.y += y
    }

    updateStartPos(){
        this.startPosition.setPosition(this.targetPosition.x, this.targetPosition.y)
    }

    interpolate(alpha: number){
        this.position.interpolateBetween(this.startPosition, this.targetPosition, alpha)
    }
}

class KinematicPlayer extends InterpolatedPlayer {
    private velocity: Vector2;

    private collisionHorizontal: Direction;
    private collisionVertical: Direction;

    constructor(abstract: AbstractPlayer){
        super(abstract)

        this.velocity = new Vector2(0, 0)

        this.collisionHorizontal = Direction.None
        this.collisionVertical = Direction.None
    }

    control(speed: number, jump: number, controls: Controls){
        if(controls.isHeld("left")){
            this.velocity.x -= speed
            controls.addTick("left")
        }

        if(controls.isHeld("right")){
            this.velocity.x += speed
            controls.addTick("right")
        }

        if(controls.isHeld("jump") && this.isCollisionInDirection(Direction.Down)){
            this.velocity.y -= jump
            controls.addTick("jump")
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

        const playerPath = new Rect({
            position: {
                x: this.targetPosition.x,
                y: this.targetPosition.y
            },
            size: {
                x: this.size.x,
                y: this.size.y,
            }
        })
        
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

        const playerPath = new Rect({
            position: {
                x: this.targetPosition.x,
                y: this.targetPosition.y
            },
            size: {
                x: this.size.x,
                y: this.size.y,
            }
        })

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

    correctDivergence(posX: number, posY: number){
        const divergenceTolerance = 30

        const distX = Math.abs(posX - this.targetPosition.x)
        if(distX >= divergenceTolerance && Math.abs(this.velocity.x) < 0.1){
            this.targetPosition.x = posX
        }

        const distY = Math.abs(posY - this.targetPosition.y)
        if(distY >= divergenceTolerance && Math.abs(this.velocity.y) < 0.1){
            this.targetPosition.y = posY
        }
    }

    getVelocity() {
        return new Vector2(this.velocity.x, this.velocity.y)
    }

    isCollisionInDirection(dir: Direction){
        if (dir == Direction.None){
            return false
        }

        return this.collisionHorizontal == dir || this.collisionVertical == dir
    }
}