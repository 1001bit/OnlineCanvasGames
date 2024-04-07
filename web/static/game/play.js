const main = $("main")
let roomID = 0
let gameID = 0
let websocket

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
    roomID = main.data("room-id")
    gameID = main.data("game-id")
    connectToWS()
})