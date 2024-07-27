class DeltaTimer {
    private lastTick: DOMHighResTimeStamp;

    constructor(){
        this.lastTick = performance.now();
    }

    getDeltaTime(){
        let now = performance.now();
        let dt = now - this.lastTick;
        this.lastTick = now;

        return dt;
    }
}