class Controls {
    // controls that are held right now (control -> bool (set))
    // using map instead of set here because golang doesn't have set implementation yet
    private heldControls: Map<string, boolean>;
    // controls bindings (key -> control)
    private bindings: Map<string, string>;

    // shows what controls were held and for how many ticks (control -> ticks)
    private heldControlsTicks: Map<string, number>

    constructor(){
        this.heldControls = new Map();
        this.bindings = new Map();

        this.heldControlsTicks = new Map()

        // on key press
        document.addEventListener("keypress", (e) => {
            // only single press
            if (e.repeat) {
                return;
            }
            // if no key in bindings
            if(!this.bindings.has(e.key)){
                return;
            }

            // get control from binding
            const control = this.bindings.get(e.key);
            if(control){
                this.heldControls.set(control, true);
            }
        });

        document.addEventListener("keyup", (e) => {
            if(!this.bindings.has(e.key)){
                return;
            }

            const control = this.bindings.get(e.key);
            if(control){
                this.heldControls.delete(control);
            }
        });
    }

    bindControl(key: string, control: string){
        this.bindings.set(key, control);
    }

    isHeld(control: string) {
        return this.heldControls.has(control)
    }

    addTick(control: string) {
        if(!this.heldControls.has(control)){
            return
        }

        const ticks = this.heldControlsTicks.get(control) 
        if(!ticks){
            this.heldControlsTicks.set(control, 1)
            return
        }

        this.heldControlsTicks.set(control, ticks+1)
    }

    resetHeldControlsTicks(serverTPS: number, clientTPS: number){
        const maxTicks = Math.ceil(clientTPS/serverTPS)
        for (const [control, ticks] of this.heldControlsTicks){
            if(ticks <= maxTicks){
                // delete controls, that didn't bypass the limit
                this.heldControlsTicks.delete(control)
                continue
            }
            
            // postpone ticks, that are beyond for the future, since can't send any more.
            this.heldControlsTicks.set(control, ticks - maxTicks)
        }
    }

    getHeldControlsTicks(){
        return this.heldControlsTicks
    }
}