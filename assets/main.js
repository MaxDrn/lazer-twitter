let username;
let userID;
let infoField;
window.onload = function () {
    let tweetButton = document.getElementById("tButton");
    let userField = document.getElementById("userField");
    let passwordField = document.getElementById("passwordField");
    let feed = document.getElementById("feed");
    let signIn = document.getElementById("signin");
    let message = "";
    let signOut = document.getElementById("signOut");
    signOut.style.display = "none";
    let tweets = [];
    let url = window.location.href;
    let prefix = url.split("//");
    let sockUrl = "ws://" + prefix[1].split("/")[0] + "/socket";
    let login = prefix[1].split("/")[0] + "/login";
    let connection = new WebSocket(sockUrl);
    let signup = document.getElementById("signup");
    let par = new URLSearchParams(document.location.search.substring(1));
    let user = par.get("username");
    let pass = par.get("password");
    window.history.pushState({}, document.title, "http://" + prefix[1].split("/")[0]);
    infoField = document.getElementById("info");
    infoField.style.visibility = "hidden";

    initListener(signIn, signOut, signup, login, tweetButton, message, userField, passwordField, connection);

    connection.onopen = function () {
        join(connection, user, pass);
    };

    setInterval(handleData(connection, feed, tweets, signOut, signIn, userField, passwordField), 1000 / 10);
};

function handleData(connection, feed, tweets, signOut, signIn, userField, passwordField) {
    let jsonData;
    connection.onmessage = function (evt) {
        if (evt.data !== null && evt.data !== "") {
            jsonData = JSON.parse(evt.data);
        }
        if (jsonData.typ === "message" && tweets.includes(jsonData.tweet.id) !== true) {
            let tweetDiv = document.createElement("div");
            let userDiv = document.createElement("div");
            let likebutton = document.createElement("button");
            let likecount = document.createElement("div");
            let date = document.createElement("div");
            let wrapper = document.createElement("div");
            date.className = "date";
            date.innerHTML = jsonData.tweet.time;
            date.id = "date";
            likecount.className = "likecount";
            likecount.innerHTML = jsonData.tweet.likes;
            likebutton.id = jsonData.tweet.id;
            likebutton.className = "likebutton";
            likebutton.innerHTML = parseInt(likecount.innerHTML) + " &#9786;";
            userDiv.className = "users";
            userDiv.innerHTML = jsonData.tweet.user + " | " + date.innerHTML;
            tweetDiv.className = "tweets";
            tweetDiv.innerHTML = jsonData.tweet.message;
            tweets.push(jsonData.tweet.id);
            likeListener(likebutton, connection, likecount);
            insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv, jsonData);
        } else if (jsonData.typ === "all") {
            for (let i = 0; i < jsonData.tweetObjects.length; i++) {
                if (tweets.includes(jsonData.tweetObjects[i].id) !== true) {
                    let tweetDiv = document.createElement("div");
                    let userDiv = document.createElement("div");
                    let likebutton = document.createElement("button");
                    let likecount = document.createElement("div");
                    let date = document.createElement("div");
                    let wrapper = document.createElement("div");
                    date.className = "date";
                    date.innerHTML = jsonData.tweetObjects[i].time;
                    date.id = "date";
                    likecount.className = "likecount";
                    likecount.innerHTML = jsonData.tweetObjects[i].likes;
                    likebutton.className = "likebutton";
                    likebutton.id = jsonData.tweetObjects[i].id;
                    likebutton.innerHTML = parseInt(likecount.innerHTML) + " &#9786;";
                    userDiv.className = "users";
                    userDiv.innerHTML = jsonData.tweetObjects[i].user + " | " + date.innerHTML;
                    tweetDiv.className = "tweets";
                    tweetDiv.innerHTML = jsonData.tweetObjects[i].message;
                    tweets.push(jsonData.tweetObjects[i].id);
                    likeListener(likebutton, connection);
                    insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv, jsonData);
                }
            }
        } else if (jsonData.typ === "liked") {
            let desButton = document.getElementById(jsonData.tweet.id);
            let newCount = jsonData.tweet.likes;
            desButton.innerHTML = parseInt(newCount) + " &#9786;";
        } else if (jsonData.typ === "loggedin") {
            username = jsonData.username;
            infoField.innerHTML = "logged in as: " + username;
            infoField.style.visibility = "visible";
            signIn.style.display = "none";
            userField.style.display = "none";
            passwordField.style.display = "none";
            signOut.style.display = "inline-block";
            console.log("logged in");
            userID = jsonData.id;
        } else if (jsonData.typ === "failedLogin") {
            bootbox.alert("wrong username and/or password");
            console.log("incorrect credentials");
        } else if (jsonData.typ === "registered") {
            console.log("successfully registered");
        } else if (jsonData.typ === "failedRegister") {
            bootbox.alert("An account with that name already exists");
            console.log("failed to register");
        } else if (jsonData.typ === "failedLike") {
            console.log("failed to like");
        } else if (jsonData.typ === "error") {
            console.log(jsonData.message);
        }
    };


}

function join(connection, user, pass) {
    let info = {
        typ: "join",
    };
    let stringjson = JSON.stringify(info);
    connection.send(stringjson);
    if (user !== "" && pass !== "" && user !== null && pass !== null) {
        signUp(user, pass, connection);
    }
}

function signUp(user, pass, connection) {
    let signup = {
        typ: "signUp",
        username: user,
        password: pass
    };
    let signUpString = JSON.stringify(signup);
    connection.send(signUpString)
}

function likeListener(likebutton, connection) {
    likebutton.addEventListener("click", function () {
        if (username !== "" && username !== null && username !== undefined) {
            let message = {
                typ: "like",
                userid: userID,
                tweetid: parseInt(this.id)
            };

            let strMessage = JSON.stringify(message);
            connection.send(strMessage);
        } else {
            bootbox.alert("you cannot like a tweet when not logged in");
        }
    });
}

function insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv) {
    userDiv.appendChild(likebutton);
    wrapper.appendChild(userDiv);
    wrapper.appendChild(tweetDiv);
    feed.prepend(wrapper);
}

function tweetListener(tweetButton, message, connection) {
    tweetButton.addEventListener("click", function (evt) {
        bootbox.prompt({
            title: "Please enter your tweet content",
            centerVertical: true,
            inputType: "textarea",
            callback: function (result) {
                if (result !== "" && result !== undefined && result !== null && username !== undefined && username !== "" && username !== null) {
                    message = result;

                    let time = moment().format('MMMM Do YYYY, h:mm:ss a');
                    let tweet = {
                        typ: "message",
                        tweet: {
                            id: 0,
                            time: time,
                            likes: 0,
                            user: username,
                            userid: userID,
                            message: message,
                        }
                    };

                    let stringTweet = JSON.stringify(tweet);
                    connection.send(stringTweet);
                } else if (result === "" || result === undefined || result === null || username === undefined || username === "" || username === null) {
                    bootbox.alert("Please make sure your tweet is not empty, also if not done already, log in to your account in order to write a tweet");
                }
            }
        });
    });
}

function signInListener(signIn, userField, passwordField, connection) {
    signIn.addEventListener("click", function (evt) {
        if (userField.value !== "" && passwordField.value !== "") {
            let login = {
                typ: "login",
                username: userField.value,
                password: passwordField.value
            };
            userField.value = "";
            passwordField.value = "";

            let stringLogin = JSON.stringify(login);
            connection.send(stringLogin);
        } else if (userField.value === "" || passwordField.value === "") {
            bootbox.alert("please enter your username and/or password");
        }
    });
}

function signOutListener(signOut, signIn, userField, passwordField) {
    signOut.addEventListener("click", function () {
        if (username !== "" && username !== null && username !== undefined) {
            username = "";
            infoField.innerHTML = "";
            infoField.style.visibility = "hidden";
            signIn.style.display = "inline-block";
            userField.style.display = "inline-block";
            passwordField.style.display = "inline-block";
            signOut.style.display = "none";
        } else {
            bootbox.alert("you are not logged in");
        }
    });
}

function signUpListener(signup, login) {
    signup.addEventListener("click", function () {
        if (username !== "" && username !== null && username !== undefined) {
            bootbox.alert("you are already logged in");
        } else {
            window.location.replace("http://" + login);
        }
    });
}

function initListener(signIn, signOut, signup, login, tweetButton, message, userField, passwordField, connection) {
    signInListener(signIn, userField, passwordField, connection);
    signOutListener(signOut, signIn, userField, passwordField);
    signUpListener(signup, login);
    tweetListener(tweetButton, message, connection);
}