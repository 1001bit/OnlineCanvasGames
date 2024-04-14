let gameID = 0
let eventSource

function connectToSSE(gameID){
    eventSource = new EventSource(`http://${document.location.host}/rt/sse/game/${gameID}`)

    eventSource.onopen = (event) => {
        console.log("sse connection open")
    }

    eventSource.onclose = (event) => {
        console.log("sse connection close")
    }

    eventSource.onmessage = (msg) => {
        handleMessage(msg)
    }
}

function handleMessage(msg){
    console.log("server said:", msg)
}

$("main").ready(() => {
    gameID = $("main").data("game-id")
    connectToSSE(gameID)
})

$("#create").click(() => {
    fetch("/api/room", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            gameid: gameID
        })
    })
    .then (response => {
        if(response.status != 200){
            response.json().then(data => $("#create").text(data.message))
            return
        }
        response.json().then(data => window.location.href = `/game/${gameID}/room/${data.roomid}`)
    })
})