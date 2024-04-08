const main = $("main")
const urlParams = new URLSearchParams(window.location.search);
let websocket

let gameID = 0
let roomID

function connectToWS(){
    websocket = new WebSocket(`ws://${document.location.host}/ws/game/${gameID}/room/${roomID}`)

    websocket.onopen = function(e) {
        console.log("connected to websocket")
    }

    websocket.onmessage = (msg) => {
        console.log(msg)
    }

    websocket.onclose = (event) => {
        console.log("closed event connection")
    }
}

main.ready(() => {
    gameID = main.data("game-id")
    roomID = urlParams.get("room")
    if (!roomID) {
        roomID = 0
    }

    connectToWS()
})