"use strict";
var InputType;
(function (InputType) {
    InputType["Login"] = "login";
    InputType["Register"] = "register";
})(InputType || (InputType = {}));
function postInput(type) {
    const inputData = {
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
        .then(response => {
        response.json().then(data => $("#info").html(data.body));
        if (response.status == 200) {
            window.location.reload();
        }
    });
}
$("#login").on("click", () => {
    postInput(InputType.Login);
});
$("#register").on("click", () => {
    postInput(InputType.Register);
});
$(() => {
    const usernameField = $("#username");
    const passwordField = $("#password");
    // auto select username field
    usernameField.trigger("focus");
    // select password field after pressing enter
    usernameField.on("keydown", e => {
        if (e.key == "Enter" && usernameField.val() != "") {
            passwordField.trigger("focus");
        }
    });
    // login after pressing enter
    passwordField.on("keydown", e => {
        if (e.key == "Enter" && passwordField.val() != "") {
            postInput(InputType.Login);
        }
    });
});
