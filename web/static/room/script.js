"use strict";
var RoomGui = /** @class */ (function () {
    function RoomGui() {
        var _this = this;
        this.navVisible = true;
        this.onclick();
        $("#show-nav").click(function () {
            _this.onclick();
        });
    }
    RoomGui.prototype.showMessage = function (text) {
        $("#message").text(text);
    };
    RoomGui.prototype.setNavBarVisibility = function (visibility) {
        this.navVisible = visibility;
        if (visibility) {
            $("#navigation").show();
            $("#gui").removeAttr("style");
            $("#show-nav").text("↑");
            return;
        }
        $("#navigation").hide();
        $("#gui").css("top", 0);
        $("#show-nav").text("↓");
    };
    RoomGui.prototype.onclick = function () {
        this.navVisible = !this.navVisible;
        this.setNavBarVisibility(this.navVisible);
    };
    return RoomGui;
}());
// using variable, so other scripts can use it
var roomGui = new RoomGui();
