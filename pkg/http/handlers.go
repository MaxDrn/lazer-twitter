package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClientTweet struct {
	Time    string `json: "time"`
	Likes   int    `json: "likes"`
	User    string `json: "user"`
	Message string `json: "message"`
	TweetID int    `json: "tweetid"`
}

type infos struct {
	Typ   string      `json: "typ"`
	Tweet ClientTweet `json: "tweet"`
}

type AllTweets struct {
	Typ string
	TweetObjects []infos 	`json: "tweetObjects"`
}

type ErrorMessage struct{
	Typ string		`json: "typ"`
	Message string 	`json: "message"`
}


var Database *sql.DB

var SocketSlice = make([]*websocket.Conn, 0)
func SocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	serverSocket, err := upgrader.Upgrade(w, r, nil)
	SocketSlice = append(SocketSlice, serverSocket)

	if err != nil {
		fmt.Println(err)
	}
	if serverSocket != nil {
		handleData(serverSocket, w)
	}
}

func handleData(socket *websocket.Conn, w http.ResponseWriter) {
	for {
		_, messageByte, err := socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		inf := infos{}
		err = json.Unmarshal(messageByte, &inf)


		if err != nil{
			panic(err.Error())
		}
		if len(inf.Typ) <= 0 || len(inf.Tweet.Message) <= 0{
			err := ErrorMessage{
				Typ: "error",
				Message: "Empty input, please check the tweet requirements",
			}
			byteErr, _ := json.Marshal(err)
			_ = socket.WriteMessage(1, byteErr)
		}
		if inf.Typ == "message" {
			insertIntoDatabase(&inf)
			for _, val := range SocketSlice{
				msg, err := json.Marshal(inf)
				if err != nil{
					panic(err.Error())
				}
				_ = val.WriteMessage(1, msg)
			}
		} else if inf.Typ == "join" {
			Tweets := make([]infos, 0)
			rows := getFromDatabase()
			for rows.Next() {
				copyTweet := infos{}
				err = rows.Scan(&copyTweet.Typ, &copyTweet.Tweet.Time, &copyTweet.Tweet.Likes, &copyTweet.Tweet.User, &copyTweet.Tweet.Message, &copyTweet.Tweet.TweetID)
				if err != nil{
					panic(err.Error())
				}
				Tweets = append(Tweets, copyTweet)
			}
			aTweets := AllTweets{Typ: "all", TweetObjects: Tweets,}

			stringATweets, _ := json.Marshal(aTweets)

			err := socket.WriteMessage(1, stringATweets)

			if err != nil{
				panic(err.Error())
			}
		} else if inf.Typ == "like"{
			id := inf.Tweet.TweetID
			likeTweet(id)
			for _, val := range SocketSlice{
				rows := getRow(id)

				for rows.Next() {
					copyTweet := infos{}
					err = rows.Scan(&copyTweet.Typ, &copyTweet.Tweet.Time, &copyTweet.Tweet.Likes, &copyTweet.Tweet.User, &copyTweet.Tweet.Message, &copyTweet.Tweet.TweetID)
					copyTweet.Typ = "liked"
					if err != nil{
						panic(err.Error())
					}

					likeTweet, _ := json.Marshal(copyTweet)
					_ = val.WriteMessage(1, likeTweet)
				}

			}
		}
	}
}

func insertIntoDatabase(inf *infos){
	rows, err := Database.Query(`
	SELECT * from Tweets;
	`)

	if err != nil{
		panic(err.Error())
	}
	id := 0
	for rows.Next(){
		id += 1
	}
	inf.Tweet.TweetID = id
	Database.Exec(`
	INSERT INTO Tweets VALUES($1, $2, $3, $4, $5, $6)
	`, inf.Typ, inf.Tweet.Time, inf.Tweet.Likes, inf.Tweet.User, inf.Tweet.Message, id)
}

func getFromDatabase() *sql.Rows {
	result, err := Database.Query(`
	SELECT * FROM Tweets;
	`)

	if err != nil {
		panic(err.Error())
		return nil
	}

	return result
}

func likeTweet(id int){
	Database.Exec(`
	UPDATE Tweets SET Likes=Likes+1 WHERE TweetID=$1;
	`, id)
}

func getRow(id int) *sql.Rows{
	result, err := Database.Query(`
		SELECT * FROM Tweets WHERE TweetID=$1;
	`, id)

	if err != nil{
		panic(err.Error())
	}

	return result
}
