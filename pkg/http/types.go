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
