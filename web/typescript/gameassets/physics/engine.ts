class PhysicsEngine {
    staticRects: Map<number, PhysicalRect>;
    kinematicRects: Map<number, KinematicRect>;
    interpolatedRects: Map<number, InterpolatedRect>;

    serverTickAccumulator: number;

    constructor(){
        this.staticRects = new Map();
        this.kinematicRects = new Map();
        this.interpolatedRects = new Map();

        this.serverTickAccumulator = 0;
    }

    insertStaticRect(id: number, rect: PhysicalRect){
        this.staticRects.set(id, rect)
    }

    insertInterpolatedRect(id: number, rect: InterpolatedRect){
        this.interpolatedRects.set(id, rect)
    }

    insertKinematicRect(id: number, rect: KinematicRect){
        this.insertInterpolatedRect(id, rect)
        this.kinematicRects.set(id, rect)
    }

    deleteRect(id: number){
        this.staticRects.delete(id)
        this.kinematicRects.delete(id)
    }

    getRect(id: number): Rect | undefined {
        const rect = this.kinematicRects.get(id)
        if(rect){
            return rect
        }
        return this.staticRects.get(id)
    }

    tick(dt: number, serverTPS: number, constants: PhysicsConstants){
        // HACK: Use fixed timestep for this loop
        for (const [_id, rect] of this.kinematicRects){
            this.applyGravityToVel(rect, constants.gravity, dt)
            this.applyFrictionToVel(rect, constants.friction)
            this.applyCollisions(rect, dt)
            this.applyVelToPos(rect, dt)
        }

        // Interpolation
        this.serverTickAccumulator += dt;
        for (const [_id, rect] of this.interpolatedRects){
            let alpha = this.serverTickAccumulator/(1000/serverTPS)
            alpha = Math.min(1, alpha)
            rect.interpolate(alpha)
        }
    }

    serverUpdate(movedRects: Map<number, PhysicalRect>, serverTPS: number){
        if(this.serverTickAccumulator >= (1000/serverTPS)){
            this.serverTickAccumulator %= (1000/serverTPS)
        }

        for(const [_id, rect] of this.interpolatedRects){
            rect.updateStartPos()
        }

        for(const [key, val] of Object.entries(movedRects)){
            const id = Number(key)
            const serverRect = val as PhysicalRect

            const staticRect = this.staticRects.get(id)
            if(staticRect){
                staticRect.setPosition(serverRect.position.x, serverRect.position.y)
                continue
            }

            const kinematicRect = this.kinematicRects.get(id)
            if(kinematicRect){
                // TODO: Correct sometimes
                const correct = false
                if(correct){
                    kinematicRect.setTargetPos(serverRect.position.x, serverRect.position.y)
                }
                continue
            }

            const interpolatedRect = this.interpolatedRects.get(id)
            if(interpolatedRect){
                interpolatedRect.setTargetPos(serverRect.position.x, serverRect.position.y)
            }
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
}