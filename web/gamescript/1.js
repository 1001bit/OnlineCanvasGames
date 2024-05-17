canvas.setBackgroundColor(RGB(50, 150, 50))

button = new RectangleShape(300, 200, false)
button.setColor(RGB(150, 20, 20))
canvas.level.insertDrawable(button, 1)

text = new Text("0 clicks", 48)
canvas.level.insertDrawable(text, 1)

// button click
canvas.canvas.addEventListener("click", e => {
    let [x, y] = canvas.getLevelMousePos()
    
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

button.setPosition((window.innerWidth - button.rect.width)/2, (window.innerHeight - button.rect.height)/2)
text.setPosition(button.rect.left + 10, button.rect.top + 10)