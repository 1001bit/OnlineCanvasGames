class SmoothCamera {
    private position: Vector2;
    private target: KinematicPlayer | undefined;

    private strength: number;

    constructor(){
        this.position = new Vector2(0, 0);
        this.target = undefined;

        this.strength = 0.9
    }

    setTarget(target: KinematicPlayer){
        this.target = target
    }

    update(dt: number){
        if(this.target == undefined){
            return
        }
        
        const centerPos = new Vector2(
            this.target.getPosition().x + this.target.getSize().x/2,
            this.target.getPosition().y + this.target.getSize().y/2
        )

        this.position.interpolateBetween(this.position, centerPos, Math.pow(this.strength, dt))
    }

    getPosition() {
        return this.position
    }
}