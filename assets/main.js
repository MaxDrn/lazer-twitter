window.onload = function () {
    let tweetButton = document.getElementById("tButton");
    let userField = document.getElementById("userField");
    let feed = document.getElementById("feed");
    let username = "";
    let message;
    let tweets = [];
    let url = window.location.href;
    let prefix = url.split("//");
    let sockUrl = "ws://" + prefix[1] + "socket";
    let connection = new WebSocket(sockUrl);
    tweetButton.addEventListener("click", function (evt) {
        bootbox.prompt({
            title: "Please enter your tweet content",
            centerVertical: true,
            inputType: "textarea",
            callback: function (result) {
                username = userField.value;
                if (result != "" && result != undefined && result != null && username != undefined && username != "" && username != null) {
                    message = result;

                    let time = moment().format('MMMM Do YYYY, h:mm:ss a');
                    let tweet = {
                        typ: "message",
                        tweet: {
                            id: 0,
                            time: time,
                            likes: 0,
                            user: username,
                            message: message,
                        }
                    };

                    let stringTweet = JSON.stringify(tweet);
                    connection.send(stringTweet);
                } else if (result == "" || result == undefined || result == null || username == undefined || username == "" || username == null) {
                    bootbox.alert("Please make sure your tweet is not empty, also if not done already, set a username in order to be able to tweet!");
                }
            }
        });
    });


    connection.onopen = function () {
        join(connection, username);
    };

    setInterval(handleData(connection, feed, tweets), 1000 / 10);
};

function handleData(connection, feed, tweets) {
    let jsonData;
    let likedTweets = [];
    connection.onmessage = function (evt) {
        jsonData = JSON.parse(evt.data);
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
            addListener(likebutton, likedTweets, connection, likecount);
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
                    addListener(likebutton, likedTweets, connection, likecount);
                    insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv, jsonData);
                }
            }
        } else if (jsonData.typ == "liked") {
            let desButton = document.getElementById(jsonData.tweet.id);
            let newCount = jsonData.tweet.likes;
            desButton.innerHTML = parseInt(newCount) + " &#9786;";
        } else if (jsonData.typ == "error"){
            console.log(jsonData.message);
        }
    };


}

function join(connection) {
    let info = {
        typ: "join",
    };
    let stringjson = JSON.stringify(info);
    connection.send(stringjson);
}

function addListener(likebutton, likedTweets, connection) {
    likebutton.addEventListener("click", function () {
        let liked = false;
        for (let i = 0; i < likedTweets.length; i++) {
            if (likedTweets[i] == likebutton.id) {
                liked = true;
                return;
            }
            if (liked == true) {
                break;
            }
        }
        if (liked != true) {
            let message = {
                typ: "like",
                tweetid: parseInt(this.id)
            };

            let strMessage = JSON.stringify(message);
            connection.send(strMessage);
            likedTweets.push(likebutton.id);
        }
    });
}

function insertInHTML(userDiv, wrapper, feed, tweets, likebutton, tweetDiv) {
    userDiv.appendChild(likebutton);
    wrapper.appendChild(userDiv);
    wrapper.appendChild(tweetDiv);
    feed.prepend(wrapper);
}