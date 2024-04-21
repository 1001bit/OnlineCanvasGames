class GameSSE {
    eventSource
    rooms

    constructor(rooms){
        this.rooms = rooms
    }

    openConnection(gameID){
        this.eventSource = new EventSource(`http://${document.location.host}/rt/sse/game/${gameID}`)
        const ES = this.eventSource

        ES.onopen = (e) => {this.handleOpen()}
        ES.onclose = (e) => {this.handleClose()}
        ES.onmessage = (e) => {this.handleMessage(e.data)}
    }

    handleOpen(){
        console.log("sse open")
    }

    handleClose(){
        console.log("sse close")
    }

    handleMessage(msg){
        msg = JSON.parse(msg)

        if (msg.type == "rooms"){
            rooms.updateRoomList(msg.body)
        } 
    }
}