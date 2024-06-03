class Game {
    constructor(gameID, roomID){
        this.websocket = new GameWebSocket(gameID, roomID)
        this.gui = new Gui()
        this.canvas = new GameCanvas("canvas")

        this.initWebsocket()
    }

    handleGameMessage = (type, body) => {} 

    initWebsocket(){
        this.websocket.handleClose = () => {
            this.stopWithText("Connection closed!")
        }

        this.websocket.handleError = () => {
            this.stopWithText("Something went wrong!")
        }

        this.websocket.handleMessage = (msg) => {
            this.handleRawMessage(msg)
        }
    }

    handleRawMessage(msg){
        if (msg.type == "close"){
            this.stopWithText(msg.body)
            return
        }
        this.handleGameMessage(msg.type, msg.body)
    }

    stopWithText(text){
        this.canvas.stop()
        this.gui.showMessage(text)
        this.gui.setNavBarVisibility(true)
    }
}