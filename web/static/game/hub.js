const main = $("main")
let gameID = 0
let eventSource

function connectToSSE(){
    eventSource = new EventSource(`http://${document.location.host}/sse/game/${gameID}/hub`)

    eventSource.onopen = (event) => {
        console.log("sse conenction open")
    }

    eventSource.onclose = (event) => {
        console.log("sse conenction close")
    }

    eventSource.onmessage = (msg) => {
        handleMessage(msg)
    }
}

function handleMessage(msg){
    console.log("server said:", msg)
}

main.ready(() => {
    gameID = main.data("game-id")
    connectToSSE()
})