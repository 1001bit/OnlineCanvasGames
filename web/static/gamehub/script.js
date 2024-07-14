"use strict";
var Rooms = /** @class */ (function () {
    function Rooms(roomListID) {
        this.roomList = $("#".concat(roomListID));
    }
    Rooms.prototype.newRoom = function (roomJSON) {
        var room = $(".sample.room").clone();
        room.removeClass("sample");
        room.find(".owner").text("".concat(roomJSON.owner, "'s room"));
        room.find(".clients").text("".concat(roomJSON.clients, "/").concat(roomJSON.limit, " players"));
        room.find(".join").attr("href", "/game/".concat(gameID, "/room/").concat(roomJSON.id));
        return room;
    };
    Rooms.prototype.updateRoomList = function (listJSON) {
        var _this = this;
        this.roomList.empty();
        listJSON.forEach(function (roomJSON) {
            _this.roomList.append(_this.newRoom(roomJSON));
        });
    };
    return Rooms;
}());
var GameSSE = /** @class */ (function () {
    function GameSSE(rooms) {
        this.rooms = rooms;
        this.eventSource = null;
    }
    GameSSE.prototype.openConnection = function (gameID) {
        var _this = this;
        var protocol = location.protocol == "https:" ? "https:" : "http:";
        this.eventSource = new EventSource("".concat(protocol, "//").concat(document.location.host, "/rt/sse/game/").concat(gameID));
        var sse = this.eventSource;
        sse.onopen = function (_e) { _this.handleOpen(); };
        sse.onerror = function (_e) { _this.handleClose(); };
        sse.onmessage = function (e) { _this.handleMessage(e.data); };
    };
    GameSSE.prototype.handleOpen = function () { };
    GameSSE.prototype.handleClose = function () { };
    GameSSE.prototype.handleMessage = function (data) {
        var msg = JSON.parse(data);
        if (msg.type == "rooms") {
            var array = msg.body;
            rooms.updateRoomList(array);
        }
    };
    return GameSSE;
}());
/// <reference path="rooms.ts"/>
/// <reference path="sse.ts"/>
var rooms = new Rooms("rooms");
var sse = new GameSSE(rooms);
var gameID = 0;
$(function () {
    gameID = $("main").data("game-id");
    sse.openConnection(gameID);
});
$("#random").on("click", joinRandomRoom);
$("#create").on("click", createRoom);
function joinRandomRoom() {
    var roomList = rooms.roomList.find(".room");
    if (roomList.length == 0) {
        $("#random").text("No rooms yet!");
        return;
    }
    var room = roomList[Math.floor(Math.random() * roomList.length)];
    if (!room) {
        return;
    }
    var href = $(room).find(".join").attr("href");
    if (!href) {
        return;
    }
    window.location.href = href;
}
function createRoom() {
    fetch("/api/game/".concat(gameID, "/room"), {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    })
        .then(function (response) {
        if (response.status != 200) {
            response.json().then(function (data) { return $("#create").text(data.body); });
            return;
        }
        response.json().then(function (data) { return window.location.href = "/game/".concat(gameID, "/room/").concat(data.body); });
    });
}
