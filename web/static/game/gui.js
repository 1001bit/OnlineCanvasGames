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

    setNavBarVisibility(visibility){
        this.navvisible = visibility

        if(visibility){
            $("#navigation").show()
            $("#gui").removeAttr("style")
            $("#show-nav").text("↑")
            return
        }
        $("#navigation").hide()
        $("#gui").css("top", 10)
        $("#show-nav").text("↓")
    }

    onclick(){
        this.navvisible = !this.navvisible
        this.setNavBarVisibility(this.navvisible) 
    }
}