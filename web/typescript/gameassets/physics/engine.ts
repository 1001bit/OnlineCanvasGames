class PhysicsEngine {
    staticRects: Map<number, PhysicalRect>;
    kinematicRects: Map<number, KinematicRect>;
    interpolatedRects: Map<number, InterpolatedRect>;

    constructor(){
        this.staticRects = new Map();
        this.kinematicRects = new Map();
        this.interpolatedRects = new Map();
    }

    insertStaticRect(id: number, rect: PhysicalRect){
        this.staticRects.set(id, rect)
    }

    insertInterpolatedRect(id: number, rect: InterpolatedRect){
        this.interpolatedRects.set(id, rect)
    }

    insertKinematicRect(id: number, rect: KinematicRect){
        this.kinematicRects.set(id, rect)
    }

    deleteRect(id: number){
        this.staticRects.delete(id)
        this.kinematicRects.delete(id)
        this.interpolatedRects.delete(id)
    }

    update(dt: number, constants: PhysicsConstants){
        for (const [_id, rect] of this.kinematicRects){
            this.applyGravityToVel(rect, constants.gravity, dt)
            this.applyFrictionToVel(rect, constants.friction)
            this.applyCollisions(rect, dt)
            this.applyVelToPos(rect, dt)
        }
    }

    updateKinematicsInterpolation(){
        for (const [_id, rect] of this.kinematicRects){
            rect.updateStartPos()
        }
    }

    updateInterpolatedInterpolation(){
        for (const [_id, rect] of this.interpolatedRects){
            rect.updateStartPos()
        }
    }

    interpolate(interpolatedAlpha: number, kinematicAlpha: number){
        for (const [_id, rect] of this.kinematicRects){
            rect.interpolate(kinematicAlpha)
        }

        for (const [_id, rect] of this.interpolatedRects){
            rect.interpolate(interpolatedAlpha)
        }
    }   

    applyGravityToVel(rect: KinematicRect, gravity: number, dt: number){
        if(!rect.doApplyGravity){
            rect
        }

        rect.velocity.y += gravity * dt
    }

    applyFrictionToVel(rect: KinematicRect, friction: number){
        if(!rect.doApplyFriction){
            return
        }

        rect.velocity.x -= rect.velocity.x * friction

        // also do friction on y axis if non gravitable
        if(!rect.doApplyGravity){
            rect.velocity.y -= rect.velocity.y * friction
        }
    }

    applyCollisions(rect: KinematicRect, dt: number){
        if(!rect.doApplyCollisions){
            return
        }

        for (const [_id, staticRect] of this.staticRects){
            collideKinematicWithStatic(rect, staticRect, dt)
        }
    }

    applyVelToPos(rect: KinematicRect, dt: number){
        const posX = rect.targetPosition.x + rect.velocity.x * dt
        const posY = rect.targetPosition.y + rect.velocity.y * dt
        rect.setTargetPos(posX, posY)
    }

    setMultiplePositions(positions: Map<number, Vector2>){
        for(const [key, val] of Object.entries(positions)){
            const id = Number(key)
            const position = val as Vector2

            const staticRect = this.staticRects.get(id)
            if(staticRect){
                staticRect.setPosition(position.x, position.y)
                continue
            }

            const kinematicRect = this.kinematicRects.get(id)
            if(kinematicRect){
                // TODO: Correct sometimes
                const correct = false
                if(correct){
                    kinematicRect.setTargetPos(position.x, position.y)
                }
                continue
            }

            const interpolatedRect = this.interpolatedRects.get(id)
            if(interpolatedRect){
                interpolatedRect.setTargetPos(position.x, position.y)
            }
        }
    }
}