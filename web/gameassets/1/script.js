class Clicker {
    // BUG: Game draws itself only after first click

    constructor(){
        const layers = 2
        this.clicks = 0

        this.game = new Game(layers)
        this.game.canvas.setBackgroundColor(RGB(60, 70, 70))

        this.websocket = new GameWebSocket()

        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)

        this.initDrawables()
    }

    initWebsocket(gameID, roomID){
        this.game.initWebsocket(this.websocket, gameID, roomID, (type, body) => {
            switch (type) {
                case "clicks":
                    this.click(body)
                    break;
            
                default:
                    break;
            }
        })
    }

    initDrawables(){
        const game = this.game

        this.button = new RectangleShape(300, 200, false)
        const button = this.button

        button.setColor(RGB(150, 150, 40))
        button.setPosition((window.innerWidth - button.rect.size.x)/2, (window.innerHeight - button.rect.size.y)/2)
        game.canvas.insertDrawable(button, 0)

        this.text = new Text("0 clicks", 48)
        const text = this.text

        text.setPosition(button.rect.position.x + 10, button.rect.position.y + 10)
        game.canvas.insertDrawable(text, 1)

        // button click
        game.canvas.canvas.addEventListener("click", e => {
            let mPos = game.canvas.getLevelMousePos()
            
            if (button.rect.containsPoint(mPos.x, mPos.y)){
                this.click(this.clicks+1)
                this.websocket.sendMessage("click", 0)
            }
        })
    }

    click(clicks){
        this.clicks = clicks
        this.text.setString(`${this.clicks} clicks`)
        this.game.canvas.draw()
    }
}

const clicker = new Clicker()