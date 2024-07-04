class Ticker {
    constructor() {
        this.timer = new DeltaTimer()
    }

    tick(callback){
        let dt = this.timer.getDeltaTime()
        callback(dt)
        requestAnimationFrame(() => this.tick(callback))
    }
}