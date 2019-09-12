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
			expectedOutput: []byte(`{"typ":"loggedin","username":"MaxDrn"}`),
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
