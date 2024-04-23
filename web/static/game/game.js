const rooms = new Rooms("rooms")
const sse = new GameSSE(rooms)
let gameID = 0

$("main").ready(() => {
    gameID = $("main").data("game-id")
    sse.openConnection(gameID)
})

function joinRandomRoom(){
    const roomList = rooms.roomList.find(".room")
    if (roomList.length == 0){
        $("#random").text("No rooms yet!")
        return
    }
    const room = roomList[Math.floor(Math.random() * roomList.length)]
    window.location.href = $(room).find(".join").attr("href")

    return

    fetch(`/api/game/${gameID}/room`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
        },
    })
    .then (response => {
        if(response.status != 200){
            response.json().then(data => $("#random").text(data.message))
            return
        }
        response.json().then(data => window.location.href = `/game/${gameID}/room/${data.roomid}`)
    })
}

function createRoom(){
    fetch(`/api/game/${gameID}/room`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    })
    .then (response => {
        if(response.status != 200){
            response.json().then(data => $("#create").text(data.message))
            return
        }
        response.json().then(data => window.location.href = `/game/${gameID}/room/${data.roomid}`)
    })
}

$("#random").click(joinRandomRoom)
$("#create").click(createRoom)