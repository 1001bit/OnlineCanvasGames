class DrawablesLayers {
    constructor(layersCount) {
        this.layers = []
        for (let i = 0; i < layersCount; i++){
            this.layers.push([])
        }
    }

    draw(ctx){
        this.layers.forEach(layer => {
            layer.forEach(drawable => {
                drawable.draw(ctx)
            })
        })
    }

    insertDrawable(drawable, layer){
        this.layers[layer].push(drawable)
    }
}