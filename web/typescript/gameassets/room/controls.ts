class Controls {
    heldControls: Map<string, boolean>;
    controlsCoeffs: Map<string, number>
    bindings: Map<string, string>;

    constructor(){
        // using map instead of set here because golang doesn't have set implementation yet
        this.heldControls = new Map();
        this.controlsCoeffs = new Map();
        this.bindings = new Map();

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

    resetCoeffs(){
        this.controlsCoeffs.clear()
    }

    updateCoeffs(serverTPS: number, clientTPS: number){
        for(const [control, _] of this.heldControls){
            let coeff = this.controlsCoeffs.get(control)
            if(coeff == undefined){
                coeff = 0
            }

            this.controlsCoeffs.set(control, coeff + serverTPS/clientTPS)
        }
    }

    getCoeffs(){
        return this.controlsCoeffs;
    }
}