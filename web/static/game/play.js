const main = $("main")
const urlParams = new URLSearchParams(window.location.search);
let websocket

let gameID = 0
let roomID

function connectToWS(){
    websocket = new WebSocket(`ws://${document.location.host}/ws/game/${gameID}/room/${roomID}`)

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
    console.log("server said:", msg)
}

main.ready(() => {
    gameID = main.data("game-id")
    roomID = urlParams.get("room")
    if (!roomID) {
        roomID = 0
    }

    connectToWS()
})