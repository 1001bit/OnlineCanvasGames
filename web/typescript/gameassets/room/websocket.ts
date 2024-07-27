class GameWebSocket {
    private active: boolean;
    private websocket: WebSocket | null;

    constructor(){
        this.active = false;
        this.websocket = null;
    }

    openConnection(gameID: number, roomID: number){
        const protocol = location.protocol == "https:" ? "wss:" : "ws:" 

        this.websocket = new WebSocket(`${protocol}//${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)
        const ws = this.websocket

        ws.onopen = _e => {
            this.active = true
        }

        ws.onclose = _e => {
            if (!this.active){
                return
            }

            this.handleClose("Connection closed")
            
            this.active = false
        }

        ws.onerror = _e => {
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

    sendMessage(type: string, body: string){
        if(!this.active){
            return
        }

        const ws = this.websocket
        if(!ws){
            return
        }

        ws.send(JSON.stringify({
            type: type,
            body: body,
        }))
    }

    handleClose = (_body: string) => {}
    handleMessage = (_type: string, _body: any) => {}
}