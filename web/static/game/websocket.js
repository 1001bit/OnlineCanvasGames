class GameWebSocket {
    constructor(){
        this.open = false
    }

    openConnection(roomID, gameID){
        const protocol = location.protocol == "https:" ? "wss:" : "ws:" 

        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)
        const ws = this.websocket

        ws.onopen = (e) => {
            this.open = true
            this.handleOpen()
        }
        ws.onclose = (e) => {
            this.open = false
            this.handleClose()
        }
        ws.onerror = (e) => {
            this.open = false
            this.handleError()
        }
        ws.onmessage = (e) => {this.handleMessage(JSON.parse(e.data))}
    }

    sendMessage(type, body){
        this.websocket.send(JSON.stringify({
            type: type,
            body: body,
        }))
    }

    isOpen(){
        return this.open
    }

    handleOpen = () => {}
    handleClose = () => {}
    handleError = () => {}
    handleMessage = (msg) => {}
}