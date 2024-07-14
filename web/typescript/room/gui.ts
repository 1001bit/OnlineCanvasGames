class RoomGui{
    navVisible: boolean;

    constructor(){
        this.navVisible = true

        this.onclick()

        $("#show-nav").click(() => {
            this.onclick()
        })
    }

    showMessage(text: string){
        $("#message").text(text)
    }

    setNavBarVisibility(visibility: boolean){
        this.navVisible = visibility

        if(visibility){
            $("#navigation").show()
            $("#gui").removeAttr("style")
            $("#show-nav").text("↑")
            return
        }
        $("#navigation").hide()
        $("#gui").css("top", 0)
        $("#show-nav").text("↓")
    }

    onclick(){
        this.navVisible = !this.navVisible
        this.setNavBarVisibility(this.navVisible) 
    }
}

// using global variable, so other scripts can use it
const roomGui = new RoomGui()