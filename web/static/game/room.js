const urlParams = new URLSearchParams(window.location.search)

function connectToWS(roomID){
    const websocket = new WebSocket(`ws://${document.location.host}/ws/room/${roomID}`)

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
    const roomID = $("main").data("room-id")

    connectToWS(roomID)
})