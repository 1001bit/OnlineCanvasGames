class Ticker {
    timer: DeltaTimer;
    
    constructor() {
        this.timer = new DeltaTimer();
    }

    tick(callback: (dt: number) => void){
        let dt = this.timer.getDeltaTime();
        callback(dt);
        requestAnimationFrame(() => this.tick(callback));
    }
}

class FixedTicker {
    accumulator: number;
    tps: number;

    constructor(tps: number){
        this.tps = tps
        this.accumulator = 0
    }

    update(dt: number, callback: (dt: number) => void){
        this.accumulator += dt
        const maxAccumulator = 1000/this.tps

        while(this.accumulator >= maxAccumulator){
            callback(maxAccumulator)
            this.accumulator -= maxAccumulator
        }
    }

    setTPS(tps: number){
        this.tps = tps
    }

    getAlpha(){
        return this.accumulator/(1000/this.tps)
    }
}