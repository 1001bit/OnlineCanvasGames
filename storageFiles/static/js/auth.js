"use strict";
const usernameField = document.getElementById("username");
const passwordField = document.getElementById("password");
// auto select username field
usernameField.focus();
// select password field after pressing enter
usernameField.addEventListener("keydown", (e) => {
    if (e.key == "Enter" && usernameField.value !== "") {
        e.preventDefault();
        passwordField.focus();
    }
});
