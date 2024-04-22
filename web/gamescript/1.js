canvas.setBackgroundColor(RGB(50, 150, 50))

rect = new RectangleShape(300, 200)
rect.setColor(RGB(150, 20, 20))
rect.setPosition((window.innerWidth - rect.rect.width)/2, (window.innerHeight - rect.rect.height)/2)
canvas.insertNewDrawable(rect)

text = new Text("0 clicks", 48)
text.setPosition(rect.rect.left + 10, rect.rect.top + 10)
canvas.insertNewDrawable(text)