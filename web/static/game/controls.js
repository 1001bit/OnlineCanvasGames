class Control {
    constructor(){
        this.isHeld = false
    }
}

class Controls {
    constructor(){
        this.controls = new Map()
        this.bindings = new Map()

        document.addEventListener("keypress", (e) => {
            if (e.repeat) {
                return
            }
            if(!this.bindings.has(e.key)){
                return
            }

            let control = this.controls.get(this.bindings.get(e.key))
            control.isHeld = true
        })

        document.addEventListener("keyup", (e) => {
            if(!this.bindings.has(e.key)){
                return
            }

            let control = this.controls.get(this.bindings.get(e.key))
            control.isHeld = false
        })
    }

    bindControl(key, control){
        this.bindings.set(key, control)
        this.controls.set(control, new Control())
    }

    isHeld(control) {
        if (!this.controls.has(control)){
            return false
        }
        return this.controls.get(control).isHeld
    }

    getControlsJSON(){
        return JSON.stringify(Object.fromEntries(this.controls.entries()))
    }
}