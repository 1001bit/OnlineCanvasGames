class GameWebSocket {
    websocket
    canvas
    connected

    constructor(canvas){
        this.connected = false
        this.canvas = canvas
    }

    openConnection(roomID, gameID){
        this.websocket = new WebSocket(`ws://${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)
        const ws = this.websocket

        ws.onopen = (e) => {this.handleOpen()}
        ws.onclose = (e) => {this.handleClose()}
        ws.onerror = (e) => {this.handleError()}
        ws.onmessage = (e) => {this.handleMessage(e.data)}
    }

    handleOpen(){
        if(this.connected){
            return
        }
        this.connected = true
        this.canvas.start()

        console.log("websocket open")
    }

    handleClose(){
        this.closeWithMessage("Connection closed!")
    }

    handleError(){
        this.closeWithMessage("Something went wrong!")
    }

    closeWithMessage(text){
        if(!this.connected){
            return
        }
        this.connected = false
        this.canvas.stop()
        
        $("#message").text(text)
        this.canvas.resize()
    }

    handleMessage(msg){
        if(!this.connected){
            return
        }

        msg = JSON.parse(msg)

        if (msg.type == "close"){
            this.closeWithMessage(msg.body)
        }
    }
}