"use strict";
var InputType;
(function (InputType) {
    InputType["Login"] = "login";
    InputType["Register"] = "register";
})(InputType || (InputType = {}));
function postInput(type) {
    var inputData = {
        username: $("#username").val(),
        password: $("#password").val(),
        type: type
    };
    fetch("/api/user", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(inputData)
    })
        .then(function (response) {
        response.json().then(function (data) { return $("#info").html(data.body); });
        if (response.status == 200) {
            window.location.reload();
        }
    });
}
$("#login").on("click", function () {
    postInput(InputType.Login);
});
$("#register").on("click", function () {
    postInput(InputType.Register);
});
