const main = $("main")
let gameID = 0
let eventSource

function connectToSSE(){
    eventSource = new EventSource(`http://${document.location.host}/sse/game/${gameID}/hub`)

    eventSource.onmessage = (msg) => {
        console.log(msg)
    }

    eventSource.onclose = (event) => {
        console.log("closed event connection")
    }
}

main.ready(() => {
    gameID = main.data("game-id")
    connectToSSE()
})