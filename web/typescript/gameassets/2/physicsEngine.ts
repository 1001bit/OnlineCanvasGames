class Physics {
    constants: PlatformerConstants;
    staticRects: Map<number, Rect>;
    kinematicRects: Map<number, KinematicRect>;

    constructor(){
        this.constants = {
            friction: 0,
            gravity: 0,
            playerSpeed: 0,
            playerJump: 0
        }

        this.staticRects = new Map()
        this.kinematicRects = new Map()
    }

    setConstants(constants: PlatformerConstants){
        this.constants = constants
    }

    insertKinematicRect(id: number, rect: KinematicRect){
        this.kinematicRects.set(id, rect)
    }

    insertStaticRect(id: number, rect: Rect){
        this.staticRects.set(id, rect)
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

    tick(dt: number){
        for (const [_id, rect] of this.kinematicRects){
            this.applyVelToPos(rect, dt)
        }
    }

    applyVelToPos(rect: KinematicRect, dt: number){
        rect.position.x += rect.velocity.x * dt
        rect.position.y += rect.velocity.y * dt
    }
}