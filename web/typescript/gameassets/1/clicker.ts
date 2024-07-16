class Clicker {
    clicks: number;

    canvas: GameCanvas;
    websocket: GameWebSocket;
    ticker: Ticker;

    drawables: Map<string, Drawable>;

    constructor(){
        const layers = 2
        
        this.clicks = 0

        this.canvas = new GameCanvas("canvas", layers)
        this.canvas.setBackgroundColor(RGB(60, 70, 70))

        this.websocket = new GameWebSocket()
        const gameID = $("main").data("game-id")
        const roomID = $("main").data("room-id")
        this.initWebsocket(gameID, roomID)

        this.drawables = new Map()
        this.initDrawables()

        this.ticker = new Ticker()
        this.ticker.tick(dt => this.tick(dt))
    }

    tick(_dt: number){
        this.canvas.draw()
    }

    initWebsocket(gameID: number, roomID: number){
        this.websocket.handleMessage = (type, body) => {
            switch (type) {
                case "clicks":
                    this.click(Number(body))
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

    stopWithText(text: string){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    initDrawables(){
        const button = new RectangleShape()
        this.drawables.set("button", button)
        button.setColor(RGB(150, 150, 40))
        button.setSize(300, 200)
        button.setPosition((window.innerWidth - button.rect.size.x)/2, (window.innerHeight - button.rect.size.y)/2)
        this.canvas.insertDrawable(button, 0, 0)

        const text = new DrawableText("0 clicks", 48)
        this.drawables.set("text", text)
        text.setPosition(button.rect.position.x + 10, button.rect.position.y + 10)
        this.canvas.insertDrawable(text, 1, 1)

        // button click
        this.canvas.canvas.addEventListener("click", _e => {
            let mPos = this.canvas.getMousePos()

            if (button.rect.containsPoint(mPos.x, mPos.y)){
                this.click(this.clicks+1)
                this.websocket.sendMessage("click", "")
            }
        })
    }

    click(clicks: number){
        this.clicks = clicks
        const text = <DrawableText> this.drawables.get("text")
        text.setString(`${this.clicks} clicks`)
    }
}

new Clicker()