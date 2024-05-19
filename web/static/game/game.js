class Game {
    constructor(){
        this.gui = new Gui()
        this.canvas = new GameCanvas("canvas")
        this.websocket = new GameWebSocket()

        this.websocket.handleClose = () => {
            this.close("Connection closed!")
        }

        this.websocket.handleError = () => {
            this.close("Something went wrong!")
        }

        this.websocket.handleMessage = (msg) => {
            this.handleMessage(msg)
        }
    }

    handleGameMessage = (type, body) => {} 

    setCanvasProperties(layers, bgColor){
        this.canvas.setBackgroundColor(bgColor)
        this.canvas.setLayersCount(layers)
    }

    getMousePos(){
        return this.canvas.getMousePos()
    }

    getLevelMousePos(){
        return this.canvas.getLevelMousePos()
    }

    insertDrawable(drawable, layer){
        this.canvas.insertDrawable(drawable, layer)
    }

    openConnection(roomID, gameID){
        this.websocket.openConnection(roomID, gameID)
    }

    handleMessage(msg){
        if (msg.type == "close"){
            this.close(msg.body)
            return
        }
        this.handleGameMessage(msg.type, msg.body)
    }

    sendMessage(type, body){
        this.websocket.sendMessage(type, body)
    }

    close(text){
        this.canvas.stop()
        this.gui.showMessage(text)
    }
}