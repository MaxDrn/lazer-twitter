window.onload = function () {
    let tweetButton = document.getElementById("tButton");
    let userButton = document.getElementById("userButton");
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
                    let date = new Date();
                    let year = date.getFullYear();
                    let month = date.getMonth();
                    let day = date.getDay();
                    let hour = date.getHours();
                    let minutes = date.getMinutes();
                    let seconds = date.getSeconds();
                    if(month <= 9 && month > 0){
                        month = "0" + date.getMonth();
                    }
                    if(day <= 9 && day > 0){
                        day = "0" + date.getDay();
                    }
                    let currentDate = year + "-" + month + "-" + day + " " + hour + ":" + minutes + ":" + seconds;
                    let tweet = {
                        typ: "message",
                        tweet: {
                            time: currentDate,
                            likes: 0,
                            user: username,
                            message: message,
                            tweetid: 0
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
        if (jsonData.Typ === "message" && tweets.includes(jsonData.Tweet.TweetID) !== true) {
            let tweetDiv = document.createElement("div");
            let userDiv = document.createElement("div");
            let likebutton = document.createElement("button");
            let likecount = document.createElement("div");
            let date = document.createElement("div");
            let wrapper = document.createElement("div");
            date.className = "date";
            date.innerHTML = jsonData.Tweet.Time;
            date.id = "date";
            likecount.className = "likecount";
            likecount.innerHTML = jsonData.Tweet.Likes;
            likebutton.id = jsonData.Tweet.TweetID;
            likebutton.className = "likebutton";
            likebutton.innerHTML = parseInt(likecount.innerHTML) + " &#9786;";
            userDiv.className = "users";
            userDiv.innerHTML = jsonData.Tweet.User + " | " + date.innerHTML;
            tweetDiv.className = "tweets";
            tweetDiv.innerHTML = jsonData.Tweet.Message;
            userDiv.appendChild(likebutton);
            wrapper.appendChild(userDiv);
            wrapper.appendChild(tweetDiv);
            feed.prepend(wrapper);
            tweets.push(jsonData.Tweet.TweetID);

            likebutton.addEventListener("click", function () {
                let liked = false;
                for(let i = 0; i < likedTweets.length; i++){
                    if(likedTweets[i] == likebutton.id){
                        liked = true;
                        return;
                    }
                    if(liked == true){
                        break;
                    }
                }
                if(liked != true){
                    let message = {
                        typ: "like",
                        tweet: {
                            time: "now",
                            likes: parseInt(likecount.innerHTML),
                            user: "",
                            message: "liked the tweet",
                            tweetid: parseInt(this.id)
                        }
                    };

                    let strMessage = JSON.stringify(message);
                    connection.send(strMessage);
                    likedTweets.push(likebutton.id);
                }
            });
        } else if (jsonData.Typ === "all") {
            for (let i = 0; i < jsonData.TweetObjects.length; i++) {
                if (tweets.includes(jsonData.TweetObjects[i].Tweet.TweetID) !== true) {
                    let tweetDiv = document.createElement("div");
                    let userDiv = document.createElement("div");
                    let likebutton = document.createElement("button");
                    let likecount = document.createElement("div");
                    let date = document.createElement("div");
                    let wrapper = document.createElement("div");
                    date.className = "date";
                    date.innerHTML = jsonData.TweetObjects[i].Tweet.Time;
                    date.id = "date";
                    likecount.className = "likecount";
                    likecount.innerHTML = jsonData.TweetObjects[i].Tweet.Likes;
                    likebutton.className = "likebutton";
                    likebutton.id = jsonData.TweetObjects[i].Tweet.TweetID;
                    likebutton.innerHTML = parseInt(likecount.innerHTML) + " &#9786;";
                    userDiv.className = "users";
                    userDiv.innerHTML = jsonData.TweetObjects[i].Tweet.User + " | " + date.innerHTML;;
                    tweetDiv.className = "tweets";
                    tweetDiv.innerHTML = jsonData.TweetObjects[i].Tweet.Message;
                    userDiv.appendChild(likebutton);
                    wrapper.appendChild(userDiv);
                    wrapper.appendChild(tweetDiv);
                    feed.prepend(wrapper);
                    tweets.push(jsonData.TweetObjects[i].Tweet.TweetID);

                    likebutton.addEventListener("click", function () {
                        let liked = false;
                        for(let i = 0; i < likedTweets.length; i++){
                            if(likedTweets[i] == likebutton.id){
                                liked = true;
                                return;
                            }
                            if(liked == true){
                                break;
                            }
                        }
                        if(liked != true){
                            let message = {
                                typ: "like",
                                tweet: {
                                    time: "now",
                                    likes: parseInt(likecount.innerHTML),
                                    user: "",
                                    message: "liked the tweet",
                                    tweetid: parseInt(this.id)
                                }
                            };

                            let strMessage = JSON.stringify(message);
                            connection.send(strMessage);
                            likedTweets.push(likebutton.id);
                        }
                    });
                }
            }
        } else if (jsonData.Typ == "liked") {
            let desButton = document.getElementById(jsonData.Tweet.TweetID);
            let newCount = jsonData.Tweet.Likes;
            desButton.innerHTML = parseInt(newCount) + " &#9786;";
        }
    };


}

function join(connection, username) {
    if (username != "" && username != null && username != undefined) {
        let info = {
            typ: "join",
            tweet: {
                time: "now",
                likes: 0,
                user: username,
                message: username + " joined",
                tweetid: -1
            }
        };
        let stringjson = JSON.stringify(info);
        connection.send(stringjson);
    } else {
        let info = {
            typ: "join",
            tweet: {
                time: "now",
                likes: 0,
                user: "",
                message: username + " joined",
                tweetid: -1
            }
        };
        let stringjson = JSON.stringify(info);
        connection.send(stringjson);
    }
}