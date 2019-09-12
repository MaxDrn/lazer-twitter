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

func (l LikeHandler) Handle(inf rawMessage) ([]byte, bool, error) {
	likeMessage := likedMessage{}
	err := json.Unmarshal(inf.Msg, &likeMessage)
	if err != nil {
		return nil, true, err
	}
	id := likeMessage.TweetId
	canLike, err := l.Database.CheckLike(id, likeMessage.UserID)
	if err != nil {
		return nil, true, err
	}

	if !canLike {
		likeFailed := likedMessage{
			Typ:     "failedLike",
			UserID:  likeMessage.UserID,
			TweetId: id,
		}
		byteLikeFail, _ := json.Marshal(likeFailed)
		return byteLikeFail, false, nil
	}

	err = l.Database.LikeTweet(id, likeMessage.UserID)
	if err != nil {
		return nil, true, errors.Wrap(err, "could not like your tweet")
	}
	tweet, err := l.Database.GetTweet(id)
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
