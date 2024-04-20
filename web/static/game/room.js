const urlParams = new URLSearchParams(window.location.search)

function connectToWS(roomID, gameID){
    const websocket = new WebSocket(`ws://${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)

    websocket.onopen = function(e) {
        console.log("ws connection open")
    }

    websocket.onclose = (event) => {
        console.log("ws connection close")
    }

    websocket.onmessage = (event) => {
        handleMessage(event.data)
    }
}

function handleMessage(message){
    msg = JSON.parse(message)
    console.log(msg)
    if (msg.type == "message"){
        $("#message").text(msg.body)
    }
}

function resizeCanvas() {
    $("#canvas").width(window.innerWidth)
    $("#canvas").height(window.innerHeight - 85)
}

$("main").ready(() => {
    const gameID = $("main").data("game-id")
    const roomID = $("main").data("room-id")

    resizeCanvas()
    connectToWS(roomID, gameID)
})

window.addEventListener('resize', resizeCanvas, false);