package http

import (
	"encoding/json"
	"lazer-twitter/persistence"

	"github.com/fid-dev/go-pflog/log"
)

func NewJoinHandler(database persistence.Database) *JoinHandler {
	return &JoinHandler{
		Database: database,
	}
}

type allTweets struct {
	Typ          string                    `json:"typ"`
	TweetObjects []persistence.ClientTweet `json:"tweetObjects"`
}

type JoinHandler struct {
	Database persistence.Database
}

func (j JoinHandler) CanHandle(raw rawMessage) bool {
	return raw.Typ == "join"
}

func (j JoinHandler) Handle(raw rawMessage) ([]byte, bool, error) {
	result, err := j.Database.GetAllTweets()
	if err != nil {
		log.Error(err.Error())
	}
	aTweets := allTweets{Typ: "all", TweetObjects: result}
	byteTweets, err := json.Marshal(aTweets)
	if err != nil {
		log.Errorf("could not convert tweets %v", err)
		return nil, false, err
	}

	return byteTweets, false, err
}
