/// <reference path="interpolatedrect.ts"/>

interface AbstractKinematicRect extends AbstractPhysicalRect {
    velocity: {
        x: number,
        y: number,
    }

    forcesToApply: {}

    collisionVertical: Direction
    collisionHorizontal: Direction
}

function isAbstractKinematicRect(obj: any) : obj is AbstractKinematicRect{
    return isAbstractPhysicalRect(obj) && "velocity" in obj
}

class KinematicRect extends InterpolatedRect {
    velocity: Vector2;

    forcesToApply: Set<ForceType>;

    collisionVertical: Direction;
    collisionHorizontal: Direction;

    constructor(abstractRect: AbstractKinematicRect){
        super(abstractRect);

        this.velocity = new Vector2(abstractRect.velocity.x, abstractRect.velocity.y);

        this.forcesToApply = new Set()
        for(const [key, _val] of Object.entries(abstractRect.forcesToApply)){
            this.forcesToApply.add(key as ForceType)
        }

        this.collisionHorizontal = Direction.None
        this.collisionVertical = Direction.None
    }

    setVelocity(x: number, y: number){
        this.velocity.setPosition(x, y)
    }

    setCollisionDir(dir: Direction){
        if(dir == Direction.Down || dir == Direction.Up){
            this.collisionVertical = dir
        } else if (dir == Direction.Left || dir == Direction.Right){
            this.collisionHorizontal = dir
        } else {
            this.collisionVertical = dir
            this.collisionHorizontal = dir
        }
    }

    doApplyForce(force: ForceType){
        return this.forcesToApply.has(force)
    }
}