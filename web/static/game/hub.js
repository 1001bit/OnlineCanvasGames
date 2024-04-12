function connectToSSE(gameID){
    let eventSource = new EventSource(`http://${document.location.host}/sse/hub/${gameID}`)

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

$("main").ready(() => {
    const gameID = $("main").data("game-id")
    connectToSSE(gameID)
})

$("#create").click(() => {
    fetch("/api/room", {
        method: "POST",
    })
    .then (response => {
        if(response.status != 200){
            return
        }
        response.json().then(data => window.location.replace(`/games/room/${data.roomid}`))
    })
})