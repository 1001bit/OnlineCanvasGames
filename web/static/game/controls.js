class Control {
    constructor(){
        this.isPressed = false
        this.isHeld = false

        this.holdPeriod = 0
    }
}

function controlValueReplacer(key, value){
    let exceptions = new Set(["isHeld", "isPressed"])
    if(exceptions.has(key)) return undefined

    return value
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
            control.isPressed = true
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

    // is held
    isHeld(control) {
        if (!this.controls.has(control)){
            return false
        }
        return this.controls.get(control).isHeld
    }

    // was pressed before server update
    isPressed(control) {
        if (!this.controls.has(control)){
            return false
        }
        return this.controls.get(control).isPressed
    }

    updateHoldTime(dt){
        this.controls.forEach((control) => {
            if(control.isHeld){
                control.holdPeriod += dt
            }
        })
    }

    updatePressedStatus(){
        this.controls.forEach((control) => {
            control.isPressed = false
        })
    }

    clear(){
        this.controls.forEach((control) => {
            control.holdPeriod = 0
        })
    }

    getControlsJSON(){
        return JSON.stringify(Object.fromEntries(this.controls.entries()), controlValueReplacer)
    }
}