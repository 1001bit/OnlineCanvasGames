class DeltaTimer {
    constructor(){
        this.lastTick = Date.now()
    }

    getDeltaTime(){
        let now = Date.now()
        let dt = now - this.lastTick
        this.lastTick = now

        return dt
    }
}