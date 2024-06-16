class Level {
    constructor(){
        this.kinematicRects = new Map()
        this.staticRects = new Map()
    }

    deleteRect(id){
        this.kinematicRects.delete(id)
        this.staticRects.delete(id)
    }

    insertStaticRect(id, rect){
        this.staticRects.set(id, rect)
    }

    insertKinematicRect(id, rect){
        this.kinematicRects.set(id, rect)
    }
}