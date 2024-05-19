class GameSSE {
    constructor(rooms){
        this.rooms = rooms
    }

    openConnection(gameID){
        const protocol = location.protocol == "https:" ? "https:" : "http:" 

        this.eventSource = new EventSource(`${protocol}//${document.location.host}/rt/sse/game/${gameID}`)
        const sse = this.eventSource

        sse.onopen = (e) => {this.handleOpen()}
        sse.onclose = (e) => {this.handleClose()}
        sse.onmessage = (e) => {this.handleMessage(e.data)}
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