class Gui{
    constructor(){
        this.navvisible = true

        this.onclick()

        $("#show-nav").click(() => {
            this.onclick()
        })
    }

    showMessage(text){
        $("#message").text(text)
    }

    onclick(){
        if(this.navvisible){
            $("#navigation").hide()
            $("#gui").css("top", 10)
            $("#show-nav").text("↓")
        } else {
            $("#navigation").show()
            $("#gui").removeAttr("style")
            $("#show-nav").text("↑")
        }

        this.navvisible = !this.navvisible

        this.resizeCanvas()
    }

    resizeCanvas = () => {}
}