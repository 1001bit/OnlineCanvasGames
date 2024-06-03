class GameWebSocket {
    constructor(gameID, roomID){
        const protocol = location.protocol == "https:" ? "wss:" : "ws:" 

        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)
        const ws = this.websocket

        ws.onopen = (e) => {}

        ws.onclose = (e) => {
            this.handleClose()
        }

        ws.onerror = (e) => {
            this.handleError()
        }

        ws.onmessage = (e) => {
            this.handleMessage(JSON.parse(e.data))
        }
    }

    sendMessage(type, body){
        this.websocket.send(JSON.stringify({
            type: type,
            body: body,
        }))
    }

    handleClose = () => {}
    handleError = () => {}
    handleMessage = (msg) => {}
}