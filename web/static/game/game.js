class Game {
    constructor(){
        this.websocket = new GameWebSocket()
        this.gui = new Gui()
        this.canvas = new GameCanvas("canvas")

        this.initWebsocket()
        this.initGui()
    }

    handleGameMessage = (type, body) => {} 

    initWebsocket(){
        this.websocket.handleClose = () => {
            this.stop("Connection closed!")
        }

        this.websocket.handleError = () => {
            this.stop("Something went wrong!")
        }

        this.websocket.handleMessage = (msg) => {
            this.handleMessage(msg)
        }
    }

    initGui(){
        this.gui.resizeCanvas = this.canvas.resize
    }

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
            this.stop(msg.body)
            return
        }
        this.handleGameMessage(msg.type, msg.body)
    }

    sendMessage(type, body){
        this.websocket.sendMessage(type, body)
    }

    stop(text){
        this.canvas.stop()
        this.gui.showMessage(text)
        this.gui.setNavBarVisibility(true)
    }
}