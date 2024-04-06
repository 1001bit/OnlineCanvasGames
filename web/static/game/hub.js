const game = $("#game")
let gameID = 0
let eventSource

function connectToSSE(){
    eventSource = new EventSource(`http://${document.location.host}/sse/game/${gameID}`)

    eventSource.onmessage = (msg) => {
        console.log(msg)
    }

    eventSource.onclose = (event) => {
        console.log("closed event connection")
    }
}

game.ready(() => {
    gameID = game.data("game-id")
    connectToSSE()
})