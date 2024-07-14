interface RoomJSON {
    owner: string;
    clients: number;
    limit: number;
    id: number;
}

class Rooms {
    roomList: JQuery<HTMLElement>;

    constructor(roomListID: string){
        this.roomList = $(`#${roomListID}`)
    }

    newRoom(roomJSON: RoomJSON) {
        const room = $(".sample.room").clone()
        room.removeClass("sample")

        room.find(".owner").text(`${roomJSON.owner}'s room`)
        room.find(".clients").text(`${roomJSON.clients}/${roomJSON.limit} players`)
        room.find(".join").attr("href", `/game/${gameID}/room/${roomJSON.id}`)

        return room
    }

    updateRoomList(listJSON: Array<RoomJSON>){
        this.roomList.empty()
        listJSON.forEach(roomJSON => {
            this.roomList.append(this.newRoom(roomJSON))
        })
    }
}