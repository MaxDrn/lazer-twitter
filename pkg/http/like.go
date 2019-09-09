package http

import (
	"encoding/json"
	"lazer-twitter/persistence"

	"github.com/pkg/errors"
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
	err = l.Database.LikeTweet(id)
	if err != nil {
		return nil, true, errors.Wrap(err, "could not like your tweet")
	}
	tweet, err := l.Database.GetRow(id)
	copyTweet := Infos{}
	if tweet != nil {
		copyTweet.Tweet = *tweet
	}
	copyTweet.Typ = "liked"
	likeTweet, err := json.Marshal(copyTweet)
	if err != nil {
		return nil, true, errors.Wrap(err, "could not convert to json")
	}

	return likeTweet, true, nil
}
