class DrawableLayers {
    layers

    constructor(layersCount) {
        this.layers = []
        for (let i = 0; i < layersCount; i++){
            this.layers.push([])
        }

        this.viewPos = new KinematicRect()
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