class GameWebSocket {
    constructor(){
        this.active = false
    }

    openConnection(gameID, roomID){
        const protocol = location.protocol == "https:" ? "wss:" : "ws:" 

        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)
        const ws = this.websocket

        ws.onopen = (e) => {
            this.active = true
        }

        ws.onclose = (e) => {
            if (!this.active){
                return
            }

            this.handleClose("Connection closed")
            
            this.active = false
        }

        ws.onerror = (e) => {
            if (!this.active){
                return
            }

            this.handleClose("Something went wrong")

            this.active = false
        }

        ws.onmessage = (e) => {
            if (!this.active) {
                return
            }

            const data = JSON.parse(e.data)

            if (data.type == "close"){
                this.handleClose(data.body)
                this.active = false
            } else {
                this.handleMessage(data.type, data.body)
            }
        }
    }

    sendMessage(type, body){
        if(!this.active){
            return
        }

        this.websocket.send(JSON.stringify({
            type: type,
            body: body,
        }))
    }

    handleClose = (body) => {}
    handleMessage = (type, body) => {}
}