package http

import (
	"encoding/json"
	"lazer-twitter/persistence"

	"github.com/fid-dev/go-pflog/log"
)

func NewLikeHandler(database persistence.Database) *LikeHandler {
	return &LikeHandler{
		Database: database,
	}
}

type LikeHandler struct {
	Database persistence.Database
}

func (l LikeHandler) CanHandle(inf rawMessage) bool {
	return inf.Typ == "like"
}

type likedMessage struct {
	Typ     string `json:"typ"`
	TweetId int    `json:"tweetid"`
}

func (l LikeHandler) Handle(inf rawMessage) ([]byte, bool, error) {
	likeMessage := likedMessage{}
	err := json.Unmarshal(inf.Msg, &likeMessage)
	id := likeMessage.TweetId
	l.Database.LikeTweet(id)
	tweet, err := l.Database.GetRow(id)
	copyTweet := Infos{}
	copyTweet.Tweet = *tweet
	copyTweet.Typ = "liked"
	likeTweet, err := json.Marshal(copyTweet)

	if err != nil {
		log.Errorf("could not convert tweet struct %v", err)
	}

	return likeTweet, true, nil
}
