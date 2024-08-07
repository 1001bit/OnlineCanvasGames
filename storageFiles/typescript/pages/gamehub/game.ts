/// <reference path="rooms.ts"/>
/// <reference path="sse.ts"/>

const rooms = new Rooms("rooms")
const sse = new GameSSE(rooms)
let gameTitle = ""

$(() => {
    gameTitle = $("main").data("game-title")
    sse.openConnection(gameTitle)
})

$("#random").on("click", joinRandomRoom)
$("#create").on("click", createRoom)

function joinRandomRoom(){
    const roomList = rooms.roomList.find(".room")
    if (roomList.length == 0){
        $("#random").text("No rooms yet!")
        return
    }

    const room = roomList[Math.floor(Math.random() * roomList.length)]
    if(!room){
        return
    }

    const href = $(room).find(".join").attr("href")
    if(!href){
        return
    }

    window.location.href = href 
}

function createRoom(){
    fetch(`/api/game/${gameTitle}/room`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    })
    .then (response => {
        if(response.status != 200){
            response.json().then(data => $("#create").text(data.body))
            return
        }
        
        response.json().then(data => window.location.href = `/game/${gameTitle}/room/${data.body}`)
    })
}