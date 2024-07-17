function isPhysicalRect(obj: any) : obj is PhysicalRect{
    return "position" in obj && "size" in obj && "doApplyCollisions" in obj
}

function isKinematicRect(obj: any) : obj is KinematicRect{
    return isPhysicalRect(obj) && "velocity" in obj
}

class PhysicalRect extends Rect {
    doApplyCollisions: boolean;

    constructor(abstractRect: PhysicalRect){
        super(abstractRect)

        this.doApplyCollisions = abstractRect.doApplyCollisions
    }
} 

class InterpolatedRect extends PhysicalRect {
    startPosition: Vector2;
    targetPosition: Vector2;

    constructor(abstractRect: PhysicalRect){
        super(abstractRect)

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

class KinematicRect extends InterpolatedRect {
    velocity: Vector2;

    doApplyGravity: boolean;
    doApplyFriction: boolean;

    constructor(abstractRect: KinematicRect){
        super(abstractRect);

        this.velocity = new Vector2(abstractRect.velocity.x, abstractRect.velocity.y);

        this.doApplyGravity = abstractRect.doApplyCollisions;
        this.doApplyFriction = abstractRect.doApplyFriction;
    }

    setVelocity(x: number, y: number){
        this.velocity.setPosition(x, y)
    }
}