package http

import "lazer-twitter/persistence"

type Infos struct {
	Typ   string                  `json:"typ"`
	Tweet persistence.ClientTweet `json:"tweet"`
}

type ErrorMessage struct {
	Typ     string `json:"typ"`
	Message string `json:"message"`
}

type likedMessage struct {
	Typ     string `json:"typ"`
	UserID  int    `json:"userid"`
	User    string `json:"user"`
	TweetId int    `json:"tweetid"`
}
