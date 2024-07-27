interface LevelConfig {
    playerSpeed: number;
    playerJump: number;
    playerGravity: number;
    playerFriction: number;
}

class Level {
    private blocks: Map<number, Block>
    private interpolatedPlayers: Map<number, InterpolatedPlayer>
    private kinematicPlayers: Map<number, KinematicPlayer>

    private config: LevelConfig;
    private playerRectID: number

    private fixedTicker: FixedTicker;
    private serverTPS: number;
    private serverAccumulator: number;

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

        this.fixedTicker = new FixedTicker(10)
        
        this.serverTPS = 0
        this.serverAccumulator = 0
    }

    setConfig(config: LevelConfig){
        this.config = config
    }

    setPlayerRectID(id: number){
        this.playerRectID = id
    }

    setTPS(client: number, server: number){
        this.serverTPS = server
        this.fixedTicker.setTPS(client)
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

    tick(dt: number, controls: Controls){
        this.serverAccumulator += dt

        // interpolate interpolated players
        const interpolatedAlpha = Math.min(this.serverAccumulator/(1000/this.serverTPS), 1)
        for (const [_, player] of this.interpolatedPlayers){
            player.interpolate(interpolatedAlpha)
        }

        // update kinematic players
        this.fixedTicker.update(dt, fixedDT => {
            for(const [rectID, player] of this.kinematicPlayers){
                // update interpolation
                player.updateStartPos()

                // Control
                if(rectID == this.playerRectID){
                    player.control(this.config.playerSpeed, this.config.playerJump, controls)
                }

                // Forces
                player.applyGravity(this.config.playerGravity, fixedDT)
                player.applyFriction(this.config.playerFriction)

                // Collisions and movement
                player.setCollisionDir(Direction.None)

                // Horizontal
                for(const [_, block] of this.blocks){
                    const dir = player.detectHorizontalCollision(block, fixedDT)
                    if(dir != Direction.None){
                        player.resolveCollision(block, dir)
                        break
                    }
                }
                player.applyHorizontalVelToPos(fixedDT)

                // Vertical
                for(const [_, block] of this.blocks){
                    const dir = player.detectVerticalCollision(block, fixedDT)
                    if(dir != Direction.None){
                        player.resolveCollision(block, dir)
                        break
                    }
                }
                player.applyVerticalVelToPos(fixedDT)
            }
        })

        // interpolate kinematic players
        const kinematicAlpha = this.fixedTicker.getAlpha()
        for (const [_, player] of this.kinematicPlayers){
            player.interpolate(kinematicAlpha)
        }
    }

    handlePlayerMovement(moved: {}, correct: boolean){
        // update interpolated rects interpolation
        this.serverAccumulator = 0
        
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
                if(correct){
                    kinematic.correctDivergence(pos.x, pos.y)
                }
            }
        } 
    }
}