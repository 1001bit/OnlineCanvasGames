let gameID = 0
let eventSource

function connectToSSE(gameID){
    eventSource = new EventSource(`http://${document.location.host}/rt/sse/game/${gameID}`)

    eventSource.onopen = (event) => {
        console.log("sse connection open")
    }

    eventSource.onclose = (event) => {
        console.log("sse connection close")
    }

    eventSource.onmessage = (msg) => {
        handleMessage(msg)
    }
}

function createRoomDiv(roomJSON){
    let room = $("<div></div>").addClass("room")
    let title = $("<h3></h3>").text(`${roomJSON.owner}'s room`)
    let link = $("<a></a>")
    let button = $("<button></button>").addClass("style-button physical small").text(`${roomJSON.clients} clients`)

    link.attr("href", `/game/${gameID}/room/${roomJSON.id}`)
    link.append(button)

    room.append(title)
    room.append(link)
    return room
}

function handleMessage(msg){
    const data = JSON.parse(msg.data)
    if (data.type == "rooms"){
        $("#rooms").empty()
        data.body.forEach(roomJSON => {
            $("#rooms").append(createRoomDiv(roomJSON))
        })
    } 
}

$("main").ready(() => {
    gameID = $("main").data("game-id")
    connectToSSE(gameID)
})

$("#random").click(() => {
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
})

$("#create").click(() => {
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
})