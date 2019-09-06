package http

import (
	"encoding/json"
	"lazer-twitter/persistence"

	"github.com/fid-dev/go-pflog/log"
	"github.com/pkg/errors"
)

var _ MessageHandler = &messageHandler{}

func NewMessageHandler(database persistence.Database) *messageHandler {
	return &messageHandler{
		database: database,
	}
}

type messageHandler struct {
	database persistence.Database
}

func (m *messageHandler) CanHandle(inf rawMessage) bool {
	return inf.Typ == "message"
}

func (m *messageHandler) Handle(inf rawMessage) ([]byte, bool, error) {
	message := Infos{}
	err := json.Unmarshal(inf.Msg, &message)
	result, err := m.database.InsertIntoDatabase(&message.Tweet)

	if err != nil {
		log.Error(err.Error())
	}
	message.Tweet.Id = result
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, false, errors.Wrapf(err, "could not convert to json %s", message)
	}
	return msg, true, nil
}
