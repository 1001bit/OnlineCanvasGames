game.setCanvasProperties(2, RGB(60, 60, 60))

button = new RectangleShape(300, 200, false)
button.setColor(RGB(150, 150, 40))
game.insertDrawable(button, 1)

text = new Text("0 clicks", 48)
game.insertDrawable(text, 1)

let clicks = 0

// button click
game.canvas.canvas.addEventListener("click", e => {
    let [x, y] = game.getLevelMousePos()
    
    if (button.rect.containsPoint(x, y)){
        clicks += 1
        text.setString(`${clicks} clicks`)
        game.sendMessage("click", 0)
    }
})

// on server message
game.handleGameMessage = (type, body) => {
    if(type == "clicks"){
        clicks = body
        text.setString(`${clicks} clicks`)
    }
}

button.setPosition((window.innerWidth - button.rect.width)/2, (window.innerHeight - button.rect.height)/2)
text.setPosition(button.rect.left + 10, button.rect.top + 10)