const rooms = new Rooms("rooms")
const sse = new GameSSE(rooms)
let gameID = 0

$("main").ready(() => {
    gameID = $("main").data("game-id")
    sse.openConnection(gameID)
})

function joinRandomRoom(){
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