const gameID = $("main").data("game-id")
const roomID = $("main").data("room-id")

const game = new Game(gameID, roomID)

game.canvas.setLayersCount(2)
game.canvas.setBackgroundColor(RGB(60, 70, 70))

button = new RectangleShape(300, 200, false)
button.setColor(RGB(150, 150, 40))
button.setPosition((window.innerWidth - button.rect.size.x)/2, (window.innerHeight - button.rect.size.y)/2)
game.canvas.drawablesLayers.insertDrawable(button, 0)

text = new Text("0 clicks", 48)
text.setPosition(button.rect.position.x + 10, button.rect.position.y + 10)
game.canvas.drawablesLayers.insertDrawable(text, 1)

let clicks = 0
function click(newClicks){
    clicks = newClicks
    text.setString(`${clicks} clicks`)
    game.canvas.draw()
}

// button click
game.canvas.canvas.addEventListener("click", e => {
    let mPos = game.canvas.getLevelMousePos()
    
    if (button.rect.containsPoint(mPos.x, mPos.y)){
        click(clicks+1)
        game.websocket.sendMessage("click", 0)
    }
})

// on server message
game.handleGameMessage = (type, body) => {
    if(type == "clicks"){
        click(body)
    }
}