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
	id := likeMessage.TweetId
	rows, errTwo := l.Database.CheckLike(id, likeMessage.Username)
	if errTwo != nil {
		return nil, true, errTwo
	}

	for rows.Next() {
		mockLike := likedMessage{}
		err := rows.Scan(&mockLike.TweetId, &mockLike.Username)
		if err != nil {
			return nil, false, err
		}
		if mockLike.TweetId == id && mockLike.Username == likeMessage.Username {
			likeFailed := likedMessage{
				Typ:      "failedLike",
				Username: likeMessage.Username,
				TweetId:  id,
			}
			byteLikeFail, _ := json.Marshal(likeFailed)
			return byteLikeFail, false, nil
		}
	}

	err = l.Database.LikeTweet(id, likeMessage.Username)
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
