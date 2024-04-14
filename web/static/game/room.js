const urlParams = new URLSearchParams(window.location.search)

function connectToWS(roomID, gameID){
    const websocket = new WebSocket(`ws://${document.location.host}/rt/ws/game/${gameID}/room/${roomID}`)

    websocket.onopen = function(e) {
        console.log("ws connection open")
    }

    websocket.onclose = (event) => {
        console.log("ws connection close")
        $("#message").text(event.reason)
    }

    websocket.onmessage = (event) => {
        handleMessage(event.data)
    }
}

function handleMessage(message){
    msg = JSON.parse(message)
    console.log("server said:", msg)
}

$("main").ready(() => {
    const gameID = $("main").data("game-id")
    const roomID = $("main").data("room-id")

    connectToWS(roomID, gameID)
})