class Controls {
    constructor(){
        // using map instead of set here because golang doesn't have set implementation yet
        this.heldControls = new Map()
        this.bindings = new Map()

        // on key press
        document.addEventListener("keypress", (e) => {
            // only single press
            if (e.repeat) {
                return
            }
            // if no key in bindings
            if(!this.bindings.has(e.key)){
                return
            }

            // get control from binding
            this.heldControls.set(this.bindings.get(e.key), true)
        })

        document.addEventListener("keyup", (e) => {
            if(!this.bindings.has(e.key)){
                return
            }

            this.heldControls.delete(this.bindings.get(e.key))
        })
    }

    bindControl(key, control){
        this.bindings.set(key, control)
    }

    isHeld(control) {
        return this.heldControls.has(control)
    }

    getHeldControls(){
        return this.heldControls
    }
}