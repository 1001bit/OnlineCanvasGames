class Clicker {
    private clicks: number;

    private canvas: GameCanvas;
    private websocket: GameWebSocket;
    private ticker: Ticker;

    private drawables: Map<string, Drawable>;

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
        this.ticker.start(dt => this.tick(dt))
    }

    private tick(_dt: number){
        this.canvas.draw()
    }

    private initWebsocket(gameID: number, roomID: number){
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

    private stopWithText(text: string){
        this.canvas.stop()
        roomGui.showMessage(text)
        roomGui.setNavBarVisibility(true)
    }

    private initDrawables(){
        const button = new RectangleShape()
        this.drawables.set("button", button)
        button.setColor(RGB(150, 150, 40))
        button.setSize(300, 200)
        button.setPosition((window.innerWidth - button.getSize().x)/2, (window.innerHeight - button.getSize().y)/2)
        this.canvas.insertDrawable(button, 0, 0)

        const text = new DrawableText("0 clicks", 48)
        this.drawables.set("text", text)
        text.setPosition(button.getPosition().x + 10, button.getPosition().y + 10)
        this.canvas.insertDrawable(text, 1, 1)

        // button click
        this.canvas.onMouseClick = (_e) => {
            let mPos = this.canvas.getMousePos()

            if (button.getRect().containsPoint(mPos.x, mPos.y)){
                this.click(this.clicks+1)
                this.websocket.sendMessage("click", "")
            }
        }
    }

    private click(clicks: number){
        this.clicks = clicks
        const text = <DrawableText> this.drawables.get("text")
        text.setString(`${this.clicks} clicks`)
    }
}

new Clicker()