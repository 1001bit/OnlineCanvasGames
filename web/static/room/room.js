const urlParams = new URLSearchParams(window.location.search)
const game = new Game(canvas)

$("main").ready(() => {
    const gameID = $("main").data("game-id")
    const roomID = $("main").data("room-id")

    game.openConnection(roomID, gameID)

    $("#navigation").hide()
    $("main").css("margin-top", 10)
})