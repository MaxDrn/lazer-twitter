package http

import (
	"encoding/json"
	"lazer-twitter/persistence"

	"github.com/fid-dev/go-pflog/log"
)

func NewLoginHandler(database persistence.Database) *LoginHandler {
	return &LoginHandler{
		database: database,
	}
}

type LoginHandler struct {
	database persistence.Database
}

func (l *LoginHandler) CanHandle(raw rawMessage) bool {
	return raw.Typ == "login" || raw.Typ == "signUp"
}

func (l *LoginHandler) Handle(raw rawMessage) ([]byte, bool, error) {
	userCredentials := persistence.Login{}
	json.Unmarshal(raw.Msg, &userCredentials)

	if raw.Typ == "login" {
		data, err := l.Login(userCredentials.Username, userCredentials.Password)
		if err != nil {
			log.Error(err.Error())
			return nil, false, err
		}
		return data, false, nil
	} else if raw.Typ == "signUp" {
		data, err := l.Register(userCredentials.Username, userCredentials.Password)
		if err != nil {
			log.Error(err.Error())
			return nil, false, err
		}
		return data, false, nil
	}

	return nil, false, nil
}
func (l *LoginHandler) Login(username string, password string) ([]byte, error) {
	data, err := l.database.LoginDatabase(username, password)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if data == true {
		loginMessage := persistence.Login{
			Typ:      "loggedin",
			Username: username,
			Password: password,
		}

		byteMsg, err := json.Marshal(loginMessage)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return byteMsg, nil
	} else if data != true {
		failedLogin := persistence.Login{
			Typ:      "failedLogin",
			Username: username,
			Password: password,
		}
		byteMsg, err := json.Marshal(failedLogin)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return byteMsg, nil
	}

	return nil, nil
}

func (l *LoginHandler) Register(username string, password string) ([]byte, error) {
	data, err := l.database.RegisterDatabase(username, password)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if data == true {
		registerMessage := persistence.Login{
			Typ:      "signedin",
			Username: username,
			Password: password,
		}

		byteMsg, err := json.Marshal(registerMessage)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return byteMsg, nil
	} else if data != true {
		failedRegister := persistence.Login{
			Typ:      "failedRegister",
			Username: username,
			Password: password,
		}
		byteMsg, err := json.Marshal(failedRegister)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		return byteMsg, nil
	}

	return nil, nil
}
