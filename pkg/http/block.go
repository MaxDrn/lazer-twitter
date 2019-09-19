package http

import (
	"encoding/json"
	"lazer-twitter/persistence"
)

func NewBlockHandler(database persistence.Database) *BlockHandler {
	return &BlockHandler{
		Database: database,
	}
}

type BlockHandler struct {
	Database persistence.Database
}

type BlockMessage struct {
	Typ        string `json:"typ"`
	ReqUserID  int    `json:"requserid"`
	UserID     int    `json:"userid"`
	BlockedIDs []int  `json:"blockedIDs"`
}

type UnblockMessage struct {
	Typ       string `json:"typ"`
	ReqUserID int    `json:"requserid"`
	UserID    int    `json:"userid"`
}

type FilteredTweets struct {
	Typ             string                    `json:"typ"`
	BlockedIDs      []int                     `json:"blockedIDs"`
	BlockedUsername string                    `json:"user"`
	CurrentBlock    int                       `json:"current"`
	Tweets          []persistence.ClientTweet `json:"tweetObjects"`
}

type UnblockedTweets struct {
	Typ    string                    `json:"typ"`
	UserID int                       `json:"userid"`
	Tweets []persistence.ClientTweet `json:"tweetObjects"`
}

func (b *BlockHandler) CanHandle(inf rawMessage) bool {
	return inf.Typ == "block" || inf.Typ == "unblock"
}

func (b *BlockHandler) Handle(inf rawMessage) ([]byte, bool, error) {
	if inf.Typ == "block" {
		blockMsg := BlockMessage{}
		err := json.Unmarshal(inf.Msg, &blockMsg)
		if err != nil {
			return nil, false, err
		}

		_, err = b.Database.InsertBlockedUser(blockMsg.ReqUserID, blockMsg.UserID)
		if err != nil {
			return nil, false, err
		}

		tweets, err := b.Database.GetTweetsForUser(blockMsg.ReqUserID)
		if err != nil {
			return nil, false, err
		}

		username, err := b.Database.UsernameFromId(blockMsg.UserID)
		if err != nil {
			return nil, false, nil
		}
		allTweets := FilteredTweets{
			Typ:             "blocked",
			BlockedIDs:      blockMsg.BlockedIDs,
			BlockedUsername: username,
			CurrentBlock:    blockMsg.UserID,
			Tweets:          tweets,
		}
		byteFilteredTweets, _ := json.Marshal(allTweets)
		return byteFilteredTweets, false, nil
	} else if inf.Typ == "unblock" {
		unblockMsg := UnblockMessage{}
		err := json.Unmarshal(inf.Msg, &unblockMsg)
		if err != nil {
			return nil, false, err
		}
		tweets, err := b.Database.GetTweetsFromUserID(unblockMsg.UserID)
		if err != nil {
			return nil, false, err
		}

		_, err = b.Database.RemoveBlockedUser(unblockMsg.ReqUserID, unblockMsg.UserID)
		if err != nil {
			return nil, false, err
		}
		unblockedTweets := UnblockedTweets{
			Typ:    "unblock",
			UserID: unblockMsg.UserID,
			Tweets: tweets,
		}
		byteUnblockedTweets, err := json.Marshal(unblockedTweets)
		if err != nil {
			return nil, false, err
		}
		return byteUnblockedTweets, false, nil
	}
	return nil, false, nil
}
