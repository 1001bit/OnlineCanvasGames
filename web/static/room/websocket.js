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
        const websocket = this.websocket

        websocket.onopen = (e) => {this.handleOpen()}
        websocket.onclose = (e) => {this.handleClose()}
        websocket.onerror = (e) => {this.handleError()}
        websocket.onmessage = (e) => {this.handleMessage(e.data)}
    }

    handleOpen(){
        if(this.connected){
            return
        }
        this.connected = true

        console.log("websocket open")
    }

    handleClose(){
        if(!this.connected){
            return
        }

        this.closeWithMessage("Connection closed!")
    }

    handleError(){
        if(!this.connected){
            return
        }

        this.closeWithMessage("Something went wrong!")
    }

    closeWithMessage(text){
        if(!this.connected){
            return
        }
        this.connected = false

        $("#message").text(text)
        this.canvas.resize()
    }

    handleMessage(msg){
        if(!this.connected){
            return
        }

        msg = JSON.parse(msg)

        if (msg.type == "message"){
            this.closeWithMessage(msg.body)
        }
    }
}