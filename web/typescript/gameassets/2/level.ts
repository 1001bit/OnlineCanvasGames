interface LevelConfig {
    playerSpeed: number;
    playerJump: number;
    playerGravity: number;
    playerFriction: number;
}

class Level {
    blocks: Map<number, Block>
    interpolatedPlayers: Map<number, InterpolatedPlayer>
    kinematicPlayers: Map<number, KinematicPlayer>

    config: LevelConfig;
    playerRectID: number

    fixedTicker: FixedTicker;
    serverAccumulator: number;

    constructor(){
        this.blocks = new Map()
        this.interpolatedPlayers = new Map()
        this.kinematicPlayers = new Map()

        this.config = {
            playerSpeed: 0,
            playerJump: 0,
            playerGravity: 0,
            playerFriction: 0,
        },
        this.playerRectID = 0

        this.fixedTicker = new FixedTicker(20)
        this.serverAccumulator = 0
    }

    setConfig(config: LevelConfig){
        this.config = config
    }

    setPlayerRectID(id: number){
        this.playerRectID = id
    }

    createPlayerRectangle(serverRect: AbstractPlayer, rectID: number){
        if(this.interpolatedPlayers.has(rectID) || this.kinematicPlayers.has(rectID)){
            return
        }

        let rectangle: RectangleShape | undefined;

        if(rectID == this.playerRectID){
            const rect = new KinematicPlayer(serverRect)

            this.kinematicPlayers.set(rectID, rect)
            rectangle = new RectangleShape(rect)
        } else {
            const rect = new InterpolatedPlayer(serverRect)

            this.interpolatedPlayers.set(rectID, rect)
            rectangle = new RectangleShape(rect)
        }

        return rectangle
    }

    disconnectPlayer(rectId: number){
        this.interpolatedPlayers.delete(rectId)
        this.kinematicPlayers.delete(rectId)
    }

    createBlockRectangle(serverRect: AbstractBlock, rectID: number){
        if(this.blocks.has(rectID)){
            return
        }

        const rect = new Block(serverRect)
        this.blocks.set(rectID, rect)

        const rectangle = new RectangleShape(rect)
        return rectangle
    }

    tick(dt: number, serverTPS: number, controls: Controls){
        this.serverAccumulator += dt

        // interpolate kinematic players
        const kinematicAlpha = this.fixedTicker.getAlpha()
        for (const [_, player] of this.kinematicPlayers){
            player.interpolate(kinematicAlpha)
        }

        // interpolate interpolated players
        const interpolatedAlpha = this.serverAccumulator/(1000/serverTPS)
        for (const [_, player] of this.interpolatedPlayers){
            player.interpolate(interpolatedAlpha)
        }

        // update kinematic players
        this.fixedTicker.update(dt, fixedDT => {
            controls.updateCoeffs(serverTPS, 1000/fixedDT)

            for(const [_, player] of this.kinematicPlayers){
                // update interpolation
                player.updateStartPos()

                // forces
                player.applyGravity(this.config.playerGravity, fixedDT)
                player.applyFriction(this.config.playerFriction)

                // control
                player.control(this.config.playerSpeed, this.config.playerJump, controls)

                // TODO: Collisions
                for(const [_, _block] of this.blocks){
                    // Here
                }
                player.setCollisionDir(Direction.None)
                
                // Move rect
                player.applyVelToPos(fixedDT)
            }
        })
    }

    handlePlayerMovement(moved: Map<number, {x: number, y: number}>){
        // update interpolated rects interpolation
        for (const [_, player] of this.interpolatedPlayers){
            player.updateStartPos()
        }
        
        // set position
        for (const [key, val] of Object.entries(moved)){
            const rectID = Number(key)
            const pos = val as {x: number, y: number}

            const interpolated = this.interpolatedPlayers.get(rectID)
            if(interpolated){
                interpolated.setTargetPos(pos.x, pos.y)
                continue
            }

            const kinematic = this.kinematicPlayers.get(rectID)
            if(kinematic){
                // TODO: Correction
                const correct = false
                if(correct){
                    kinematic.setTargetPos(pos.x, pos.y)
                }
            }
        } 
    }
}