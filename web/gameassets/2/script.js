const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")

const game = new Game(gameID, roomID)

game.canvas.setLayersCount(2)
game.canvas.setBackgroundColor(RGB(30, 200, 200))

// on server message
game.handleGameMessage = (type, body) => {
    console.log(type, body)
}