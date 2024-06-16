class GameWebSocket {
    constructor(){}

    openConnection(gameID, roomID){
        const protocol = location.protocol == "https:" ? "wss:" : "ws:" 

        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)
        const ws = this.websocket

        ws.onopen = (e) => {}

        ws.onclose = (e) => {
            this.handleClose("Connection closed")
        }

        ws.onerror = (e) => {
            this.handleClose("Something went wrong")
        }

        ws.onmessage = (e) => {
            const CLOSE = "close"
            const data = JSON.parse(e.data)
            data.type == CLOSE ? this.handleClose(data.body) : this.handleMessage(data.type, data.body)
        }
    }

    sendMessage(type, body){
        if(this.websocket.readyState !== WebSocket.OPEN){
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