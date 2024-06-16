class Game {
    constructor(layers){
        this.active = true
        this.gui = new Gui()
        this.canvas = new GameCanvas("canvas", layers)
    }

    initWebsocket(websocket, gameID, roomID, callback){
        websocket.handleClose = (body) => {
            this.stopWithText(body)
        }

        websocket.handleMessage = callback

        websocket.openConnection(gameID, roomID)
    }

    stopWithText(text){
        if(!this.active){
            return
        }
        this.active = false

        this.canvas.stop()
        this.gui.showMessage(text)
        this.gui.setNavBarVisibility(true)
    }
}