canvas.setBackgroundColor(RGB(50, 150, 50))

button = new RectangleShape(300, 200)
button.setColor(RGB(150, 20, 20))
button.setPosition((window.innerWidth - button.rect.width)/2, (window.innerHeight - button.rect.height)/2)
canvas.insertNewDrawable(button)

text = new Text("0 clicks", 48)
text.setPosition(button.rect.left + 10, button.rect.top + 10)
canvas.insertNewDrawable(text)

// button click
canvas.canvas.addEventListener("click", e => {
    let {x, y} = canvas.getMousePosition(e)
    
    if (button.rect.containsPoint(x, y)){
        console.log(1)
    }
})