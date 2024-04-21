for(let i = 0; i < 10; i++){
    rect = new RectangleShape(50, 50)
    rect.setPosition(i * 50, 0)
    rect.setColor(RGB(255 * (i/10), 0, 140 * (i/10)))
    
    canvas.insertNewDrawable(rect)
}

canvas.setDrawRate(30)