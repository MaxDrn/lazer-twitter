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
	}{
		{
			name: "Default test",
			msg: rawMessage{
				Typ: "login",
				Msg: []byte(`{"typ":"login","username":"MaxDrn","password":"Hallo"}`),
			},
			expectedOutput: []byte(`{"id":0,"typ":"loggedin","username":"MaxDrn"}`),
			expectedError:  false,
		},
	}

	m := mockDB{}
	testObj := NewUserHandler(&m)
	for _, val := range cases {
		t.Run(val.name, func(tt *testing.T) {
			tweets, _, err := testObj.Handle(val.msg)
			assert.EqualValues(tt, val.expectedOutput, tweets, "Output not as expected")
			assert.EqualValues(tt, val.expectedError, err != nil, "Error not as expected")
		})
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
		Typ: "like",
		Msg: []byte(`{"typ":"like"}`),
	}

	dataTwo := testObj.CanHandle(testMsg2)
	assert.EqualValues(t, false, dataTwo, "output not as expected")
}
