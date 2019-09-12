package http

import (
	"lazer-twitter/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ persistence.Database = &mockDB{}

func Test_UserHandler(t *testing.T) {

	cases := []struct {
		name           string
		msg            rawMessage
		expectedOutput []byte
		expectedError  bool
		returnError    bool
	}{
		{
			name: "successful login",
			msg: rawMessage{
				Typ: "login",
				Msg: []byte(`{"typ":"login","username":"MaxDrn","password":"Hallo"}`),
			},
			expectedOutput: []byte(`{"id":0,"typ":"loggedin","username":"MaxDrn"}`),
			expectedError:  false,
			returnError:    false,
		},
		{
			name: "successful Register",
			msg: rawMessage{
				Typ: "signUp",
				Msg: []byte(`{"typ":"login","username":"MaxDrn","password":"Hallo"}`),
			},
			expectedOutput: []byte(`{"id":0,"typ":"registered","username":"MaxDrn"}`),
			expectedError:  false,
			returnError:    false,
		},
		{
			name: "failed login",
			msg: rawMessage{
				Typ: "login",
				Msg: []byte(`{"typ":"login","username":"Test","password":"test"}`),
			},
			expectedOutput: []byte(`{"id":0,"typ":"failedLogin","username":"Test"}`),
			expectedError:  false,
			returnError:    true,
		},
		{
			name: "failed register",
			msg: rawMessage{
				Typ: "signUp",
				Msg: []byte(`{"typ":"signUp","username":"Test","password":"test"}`),
			},
			expectedOutput: []byte(`{"id":0,"typ":"failedRegister","username":"Test"}`),
			expectedError:  false,
			returnError:    true,
		},
	}

	m := mockDB{}
	testObj := NewUserHandler(&m)
	for _, val := range cases {
		t.Run(val.name, func(tt *testing.T) {
			if val.returnError == true {
				m.returnError = true
			}
			tweets, _, err := testObj.Handle(val.msg)
			assert.EqualValues(tt, val.expectedOutput, tweets, "Output not as expected")
			assert.EqualValues(tt, val.expectedError, err != nil, "Error not as expected")
		})
		m.returnError = false
	}
}

func Test_CanHandle_User(t *testing.T) {
	m := mockDB{}
	testObj := NewUserHandler(&m)

	testMsg := rawMessage{
		Typ: "login",
		Msg: []byte(`{"typ":"login"}`),
	}
	data := testObj.CanHandle(testMsg)
	assert.EqualValues(t, true, data, "output not as expected")

	testMsg2 := rawMessage{
		Typ: "signUp",
		Msg: []byte(`{"typ":"signUp"}`),
	}

	dataTwo := testObj.CanHandle(testMsg2)
	assert.EqualValues(t, true, dataTwo, "output not as expected")

	testMsg3 := rawMessage{
		Typ: "like",
		Msg: []byte(`{"typ":"like"}`),
	}

	dataThree := testObj.CanHandle(testMsg3)
	assert.EqualValues(t, false, dataThree, "output not as expected")
}
