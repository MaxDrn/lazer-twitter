window.onload = function () {
    let button = document.getElementById("submit");
    let username = document.getElementById("username");
    let password = document.getElementById("password");
    let url = window.location.href;
    let prefix = url.split("//");
    let home = prefix[1].split("/")[0];
    button.addEventListener("click", function () {
        if (username.value != "" && password.value != "") {
            window.location.replace("http://" + home + "?username=" + username.value + "&" + "password=" + password.value);
        }
    })
};