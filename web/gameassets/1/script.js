class Clicker {
    constructor(){
        const layers = 2
        
        this.clicks = 0

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(60, 70, 70))

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)

        this.initDrawables()

        this.ticker = new Ticker()
        this.ticker.tick(dt => this.tick(dt))
    }

    tick(dt){
        this.canvas.draw()
    }

    initWebsocket(gameID, roomID){
        this.websocket.handleMessage = (type, body) => {
            switch (type) {
                case "clicks":
                    this.click(body)
                    break;
            
                default:
                    break;
            }
        }

        this.websocket.handleClose = (body) => {
            this.stopWithText(body)
        }

        this.websocket.openConnection(gameID, roomID)
    }

    stopWithText(text){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    initDrawables(){
        this.button = new RectangleShape(300, 200, false)
        const button = this.button

        button.setColor(RGB(150, 150, 40))
        button.setPosition((window.innerWidth - button.rect.size.x)/2, (window.innerHeight - button.rect.size.y)/2)
        this.canvas.insertDrawable(button, 0)

        this.text = new Text("0 clicks", 48)
        const text = this.text

        text.setPosition(button.rect.position.x + 10, button.rect.position.y + 10)
        this.canvas.insertDrawable(text, 1)

        // button click
        this.canvas.canvas.addEventListener("click", e => {
            let mPos = this.canvas.getMousePos()

            if (button.rect.containsPoint(mPos.x, mPos.y)){
                this.click(this.clicks+1)
                this.websocket.sendMessage("click", 0)
            }
        })
    }

    click(clicks){
        this.clicks = clicks
        this.text.setString(`${this.clicks} clicks`)
    }
}

new Clicker()