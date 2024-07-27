class Ticker {
    previousTime: number;
    
    constructor() {
        this.previousTime = 0
    }

    start(callback: (dt: number) => void){
        requestAnimationFrame((time) => {
            this.tick(callback, time)
        })
    }

    tick(callback: (dt: number) => void, time: number){
        const dt = time - this.previousTime
        this.previousTime = time

        callback(dt)

        requestAnimationFrame((time) => {
            this.tick(callback, time)
        })
    }
}

class FixedTicker {
    accumulator: number;
    tps: number;

    constructor(tps: number){
        this.tps = tps
        this.accumulator = 0
    }

    update(dt: number, callback: (fixedDT: number) => void){
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