function getRandomInt(max) {
    return Math.floor(Math.random() * max);
}

canvas.setBackgroundColor(RGB(50, 150, 50))

button = new RectangleShape(300, 200)
button.setColor(RGB(150, 20, 20))
canvas.insertNewDrawable(button)

text = new Text("0 clicks", 48)
canvas.insertNewDrawable(text)

// button click
canvas.canvas.addEventListener("click", e => {
    let [x, y] = canvas.getMousePos()
    
    if (button.rect.containsPoint(x, y)){
        websocket.sendMessage("click", 0)
    }
})

// on server message
websocket.handleGameMessage = function(msg){
    if(msg.type == "clicks"){
        text.setString(`${msg.body} clicks`)
    }
}

rect = new RectangleShape(100, 100, true)
canvas.insertNewDrawable(rect)
// mouse move
canvas.gameUpdate = () => {
    let [x, y] = canvas.getMousePos()
    rect.rect.setCurrentPosition(x, y, false)
    button.setPosition((window.innerWidth - button.rect.width)/2, (window.innerHeight - button.rect.height)/2)
    text.setPosition(button.rect.left + 10, button.rect.top + 10)
}