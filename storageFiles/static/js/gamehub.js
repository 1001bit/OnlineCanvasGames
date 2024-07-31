"use strict";
class Rooms {
    constructor(roomListID) {
        this.roomList = $(`#${roomListID}`);
    }
    newRoom(roomJSON) {
        const room = $(".sample.room").clone();
        room.removeClass("sample");
        room.find(".owner").text(`${roomJSON.owner}'s room`);
        room.find(".clients").text(`${roomJSON.clients}/${roomJSON.limit} players`);
        room.find(".join").attr("href", `/game/${gameTitle}/room/${roomJSON.id}`);
        return room;
    }
    updateRoomList(listJSON) {
        this.roomList.empty();
        listJSON.forEach(roomJSON => {
            this.roomList.append(this.newRoom(roomJSON));
        });
    }
}
class GameSSE {
    constructor(rooms) {
        this.rooms = rooms;
        this.eventSource = null;
    }
    openConnection(gameTitle) {
        const protocol = location.protocol == "https:" ? "https:" : "http:";
        this.eventSource = new EventSource(`${protocol}//${document.location.host}/rt/sse/game/${gameTitle}`);
        const sse = this.eventSource;
        sse.onopen = (_e) => { this.handleOpen(); };
        sse.onerror = (_e) => { this.handleClose(); };
        sse.onmessage = (e) => { this.handleMessage(e.data); };
    }
    handleOpen() { }
    handleClose() { }
    handleMessage(data) {
        const msg = JSON.parse(data);
        if (msg.type == "rooms") {
            const array = msg.body;
            rooms.updateRoomList(array);
        }
    }
}
/// <reference path="rooms.ts"/>
/// <reference path="sse.ts"/>
const rooms = new Rooms("rooms");
const sse = new GameSSE(rooms);
let gameTitle = "";
$(() => {
    gameTitle = $("main").data("game-title");
    sse.openConnection(gameTitle);
});
$("#random").on("click", joinRandomRoom);
$("#create").on("click", createRoom);
function joinRandomRoom() {
    const roomList = rooms.roomList.find(".room");
    if (roomList.length == 0) {
        $("#random").text("No rooms yet!");
        return;
    }
    const room = roomList[Math.floor(Math.random() * roomList.length)];
    if (!room) {
        return;
    }
    const href = $(room).find(".join").attr("href");
    if (!href) {
        return;
    }
    window.location.href = href;
}
function createRoom() {
    fetch(`/api/game/${gameTitle}/room`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then(response => {
        if (response.status != 200) {
            response.json().then(data => $("#create").text(data.body));
            return;
        }
        response.json().then(data => window.location.href = `/game/${gameTitle}/room/${data.body}`);
    });
}
