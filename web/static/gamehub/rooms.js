class Rooms {
    constructor(roomListID){
        this.roomList = $(`#${roomListID}`)
    }

    newRoom(roomJSON) {
        const room = $(".sample.room").clone()
        room.removeClass("sample")

        room.find(".owner").text(`${roomJSON.owner}'s room`)
        room.find(".clients").text(`${roomJSON.clients}/${roomJSON.limit} players`)
        room.find(".join").attr("href", `/game/${gameID}/room/${roomJSON.id}`)

        return room
    }

    updateRoomList(listJSON){
        this.roomList.empty()
        listJSON.forEach(roomJSON => {
            this.roomList.append(this.newRoom(roomJSON))
        })
    }
}