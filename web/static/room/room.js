const urlParams = new URLSearchParams(window.location.search)
const canvas = new GameCanvas("canvas")
const websocket = new GameWebSocket(canvas)

$("main").ready(() => {
    const gameID = $("main").data("game-id")
    const roomID = $("main").data("room-id")

    canvas.resize()
    websocket.openConnection(roomID, gameID)
})