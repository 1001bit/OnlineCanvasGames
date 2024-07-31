class GameSSE {
    rooms: Rooms;
    eventSource: EventSource | null;

    constructor(rooms: Rooms){
        this.rooms = rooms;
        this.eventSource = null;
    }

    openConnection(gameTitle: string){
        const protocol = location.protocol == "https:" ? "https:" : "http:" 

        this.eventSource = new EventSource(`${protocol}//${document.location.host}/rt/sse/game/${gameTitle}`)
        const sse = this.eventSource

        sse.onopen = (_e) => {this.handleOpen()}
        sse.onerror = (_e) => {this.handleClose()}
        sse.onmessage = (e) => {this.handleMessage(e.data)}
    }

    handleOpen(){}
    handleClose(){}

    handleMessage(data: any){
        const msg = JSON.parse(data)

        if (msg.type == "rooms"){
            const array: Array<RoomJSON> = msg.body
            rooms.updateRoomList(array)
        } 
    }
}