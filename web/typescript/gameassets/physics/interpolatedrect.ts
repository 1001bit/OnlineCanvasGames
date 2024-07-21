/// <reference path="physicalrect.ts"/>

class InterpolatedRect extends PhysicalRect {
    startPosition: Vector2;
    targetPosition: Vector2;

    constructor(abstractRect: AbstractPhysicalRect){
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