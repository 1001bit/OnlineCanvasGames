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