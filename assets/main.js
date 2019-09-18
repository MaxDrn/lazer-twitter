let username;
let userID;
let infoField;
let blockedList;
let blockedUser = [];
let tweetElements = [];
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
    blockedList = document.getElementById("blockList");

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
        if (jsonData.typ === "message" && tweets.includes(jsonData.tweet.id) !== true && !blockedUser.includes(jsonData.tweet.userid)) {
            let tweetDiv = document.createElement("div");
            let userDiv = document.createElement("div");
            let likebutton = document.createElement("button");
            let blockbutton = document.createElement("button");
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
            blockbutton.id = jsonData.tweet.id;
            blockbutton.className = "blockbutton";
            blockbutton.innerHTML = "&#10006;";
            userDiv.className = "users";
            userDiv.innerHTML = jsonData.tweet.user + " | " + date.innerHTML;
            tweetDiv.className = "tweets";
            tweetDiv.innerHTML = jsonData.tweet.message;
            tweets.push(jsonData.tweet.id);
            likeListener(likebutton, connection, likecount);
            blockListener(blockbutton, connection, jsonData.tweet.userid);
            insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv, blockbutton);
        } else if (jsonData.typ === "all") {
            tweets = [];
            jsonData.tweetObjects.sort((a, b) => (a.time > b.time) ? 1 : -1);
            for (let i = 0; i < jsonData.tweetObjects.length; i++) {
                tweetElements = [];
                if (tweets.includes(jsonData.tweetObjects[i].id) !== true) {
                    receiveTweets(jsonData, connection, feed, tweets, i, true);
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
            console.log("info: logged in as: " + username);
            userID = jsonData.id;
            feed.innerHTML = "";
            tweets = [];
            tweetElements = [];
            jsonData.tweetObjects.sort((a, b) => (a.time > b.time) ? 1 : -1);
            for (let i = 0; i < jsonData.tweetObjects.length; i++) {
                receiveTweets(jsonData, connection, feed, tweets, i, true);
            }

            for (let i = 0; i < jsonData.blockedids.length; i++){
                let blockedUserDiv = document.createElement("div");
                let unblockButton = document.createElement("button");
                let blockedWrapper = document.createElement("div");
                blockedUserDiv.className = "users";
                blockedUserDiv.innerHTML = jsonData.blockedusernames[i];
                unblockButton.className = "likebutton";
                unblockButton.innerHTML = "unblock";
                blockedWrapper.id = "unblock" + jsonData.blockedids[i].toString();
                blockedWrapper.appendChild(blockedUserDiv);
                blockedUserDiv.appendChild(unblockButton);
                blockedList.appendChild(blockedWrapper);
                let id = blockedWrapper.id.split("unblock")[1];
                unblockListener(unblockButton, connection, parseInt(id));
            }
        } else if (jsonData.typ === "failedLogin") {
            bootbox.alert("info: wrong username and/or password");
            console.log("error: incorrect credentials");
        } else if (jsonData.typ === "registered") {
            console.log("info: successfully registered");
        } else if (jsonData.typ === "failedRegister") {
            bootbox.alert("info: An account with that name already exists");
            console.log("error: failed to register");
        } else if (jsonData.typ === "failedLike") {
            console.log("error: failed to like");
        } else if (jsonData.typ === "blocked") {
            feed.innerHTML = "";
            tweets = [];
            tweetElements = [];
            jsonData.tweetObjects.sort((a, b) => (a.time > b.time) ? 1 : -1);
            for (let i = 0; i < jsonData.tweetObjects.length; i++) {
                receiveTweets(jsonData, connection, feed, tweets, i, true);
            }
            let blockedUserDiv = document.createElement("div");
            let unblockButton = document.createElement("button");
            let blockedWrapper = document.createElement("div");
            blockedUserDiv.className = "users";
            blockedUserDiv.innerHTML = jsonData.user;
            unblockButton.className = "likebutton";
            unblockButton.innerHTML = "unblock";
            blockedWrapper.id = "unblock" + jsonData.current.toString();
            blockedWrapper.appendChild(blockedUserDiv);
            blockedUserDiv.appendChild(unblockButton);
            blockedList.appendChild(blockedWrapper);
            let id = blockedWrapper.id.split("unblock")[1];
            unblockListener(unblockButton, connection, parseInt(id));
        } else if (jsonData.typ === "unblock") {
            feed.innerHTML = "";
            jsonData.tweetObjects.sort((a, b) => (a.time > b.time) ? 1 : -1);
            for (let i = 0; i < jsonData.tweetObjects.length; i++) {
                receiveTweets(jsonData, connection, feed, tweets, i, false);
            }
            tweetElements.sort((a, b) => (a.time > b.time) ? 1 : -1);
            render(tweetElements, feed, connection, tweets);
            let unblockButton = document.getElementById("unblock" + jsonData.userid.toString());
            blockedList.removeChild(unblockButton);
        } else if (jsonData.typ === "error") {
            console.log(jsonData.message);
        }
    };
}


function receiveTweets(jsonData, connection, feed, tweets, i, renderDiv) {
    let tweetDiv = document.createElement("div");
    let userDiv = document.createElement("div");
    let likebutton = document.createElement("button");
    let blockbutton = document.createElement("button");
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
    blockbutton.id = jsonData.tweetObjects[i].userid.toString();
    blockbutton.className = "blockbutton";
    blockbutton.innerHTML = "&#10006;";
    userDiv.className = "users";
    userDiv.innerHTML = jsonData.tweetObjects[i].user + " | " + date.innerHTML;
    tweetDiv.className = "tweets";
    tweetDiv.innerHTML = jsonData.tweetObjects[i].message;
    tweets.push(jsonData.tweetObjects[i].id);
    likeListener(likebutton, connection);
    blockListener(blockbutton, connection, jsonData.tweetObjects[i].userid);
    tweetElements.push(jsonData.tweetObjects[i]);
    if (renderDiv) {
        insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv, blockbutton);
    } else {
        userDiv.appendChild(likebutton);
        userDiv.appendChild(blockbutton);
        wrapper.appendChild(userDiv);
        wrapper.appendChild(tweetDiv);
    }
}

function render(ltweets, feed, connection, tweets){
    for (let i = 0; i < tweets.length; i++){
        let tweetDiv = document.createElement("div");
        let userDiv = document.createElement("div");
        let likebutton = document.createElement("button");
        let blockbutton = document.createElement("button");
        let likecount = document.createElement("div");
        let date = document.createElement("div");
        let wrapper = document.createElement("div");
        date.className = "date";
        date.innerHTML = ltweets[i].time;
        date.id = "date";
        likecount.className = "likecount";
        likecount.innerHTML = ltweets[i].likes;
        likebutton.className = "likebutton";
        likebutton.id = ltweets[i].id;
        likebutton.innerHTML = parseInt(likecount.innerHTML) + " &#9786;";
        blockbutton.id = ltweets[i].userid.toString();
        blockbutton.className = "blockbutton";
        blockbutton.innerHTML = "&#10006;";
        userDiv.className = "users";
        userDiv.innerHTML = ltweets[i].user + " | " + date.innerHTML;
        tweetDiv.className = "tweets";
        tweetDiv.innerHTML = ltweets[i].message;
        ltweets.push(tweets[i].id);
        likeListener(likebutton, connection);
        blockListener(blockbutton, connection, ltweets[i].userid);
        insertInHTML(userDiv, wrapper, feed, ltweets, likebutton, tweetDiv, blockbutton);
    }
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

function blockListener(blockbutton, connection, id) {
    blockbutton.addEventListener("click", function () {
        if (username !== "" && username !== undefined && username !== null && id !== userID) {
            blockedUser.push(id);
            let message = {
                typ: "block",
                requserid: userID,
                userid: id,
                blockedIDs: blockedUser
            };
            let strMsg = JSON.stringify(message);
            connection.send(strMsg);
        } else if (username === "" || username === undefined || username === null) {
            bootbox.alert("info: please make sure you are logged in");
        } else if (userID === id) {
            bootbox.alert("info: you cannot block yourself");
        }
    });
}

function unblockListener(unblockbutton, connection, id) {
    unblockbutton.addEventListener("click", function (evt) {
        if (username !== "" && username !== undefined && username !== null) {
            let uid = blockedUser.indexOf(id);
            blockedUser.splice(uid, 1);
            let message = {
                typ: "unblock",
                requserid: userID,
                userid: id,
            };
            let strMsg = JSON.stringify(message);
            connection.send(strMsg);
        } else if (username === "" || username === undefined || username === null) {
            bootbox.alert("info: please make sure you are logged in");
        }
    });
}

function insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv, blockbutton) {
    userDiv.appendChild(likebutton);
    userDiv.appendChild(blockbutton);
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
                            userid: userID,
                            user: username,
                            message: message,
                        }
                    };

                    let stringTweet = JSON.stringify(tweet);
                    connection.send(stringTweet);
                } else if (result === "" || result === undefined || result === null || username === undefined || username === "" || username === null) {
                    bootbox.alert("info: Please make sure your tweet is not empty, also if not done already, log in to your account in order to write a tweet");
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
            bootbox.alert("info: please enter your username and/or password");
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
            window.location.reload();
        } else {
            bootbox.alert("info: you are not logged in");
        }
    });
}

function signUpListener(signup, login) {
    signup.addEventListener("click", function () {
        if (username !== "" && username !== null && username !== undefined) {
            bootbox.alert("info: you are already logged in");
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