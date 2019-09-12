package http

import (
	"encoding/json"
	"lazer-twitter/persistence"

	"github.com/fid-dev/go-pflog/log"
)

func NewUserHandler(database persistence.Database) *UserHandler {
	return &UserHandler{
		database: database,
	}
}

type UserHandler struct {
	database persistence.Database
}

func (l *UserHandler) CanHandle(raw rawMessage) bool {
	return raw.Typ == "login" || raw.Typ == "signUp"
}

func (l *UserHandler) Handle(raw rawMessage) ([]byte, bool, error) {
	userCredentials := persistence.Login{}
	err := json.Unmarshal(raw.Msg, &userCredentials)
	if err != nil {
		log.Error(err.Error())
	}
	if raw.Typ == "login" {
		data, err := l.Login(userCredentials.Username, userCredentials.Password)
		if err != nil {
			return nil, false, err
		}
		return data, false, nil
	} else if raw.Typ == "signUp" {
		data, err := l.Register(userCredentials.Username, userCredentials.Password)
		if err != nil {
			return nil, false, err
		}
		return data, false, nil
	}

	return nil, false, nil
}
func (l *UserHandler) Login(username string, password string) ([]byte, error) {
	id, data, err := l.database.Login(username, password)
	if err != nil {
		return nil, err
	}

	if data == true {
		loginMessage := persistence.Login{
			Uid:      id,
			Typ:      "loggedin",
			Username: username,
		}

		byteMsg, err := json.Marshal(loginMessage)
		if err != nil {
			return nil, err
		}
		return byteMsg, nil
	} else if data != true {
		failedLogin := persistence.Login{
			Uid:      id,
			Typ:      "failedLogin",
			Username: username,
		}
		byteMsg, err := json.Marshal(failedLogin)
		if err != nil {
			return nil, err
		}
		return byteMsg, nil
	}

	return nil, nil
}

func (l *UserHandler) Register(username string, password string) ([]byte, error) {
	data, err := l.database.Register(username, password)
	if err != nil {
		return nil, err
	}

	if data == true {
		registerMessage := persistence.Login{
			Typ:      "registered",
			Username: username,
		}

		byteMsg, err := json.Marshal(registerMessage)
		if err != nil {
			return nil, err
		}
		return byteMsg, nil
	} else if data != true {
		failedRegister := persistence.Login{
			Typ:      "failedRegister",
			Username: username,
		}
		byteMsg, err := json.Marshal(failedRegister)
		if err != nil {
			return nil, err
		}
		return byteMsg, nil
	}

	return nil, nil
}
