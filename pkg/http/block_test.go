package http

import (
	"lazer-twitter/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ persistence.Database = &mockDB{}

func Test_BlockHandler(t *testing.T) {

	cases := []struct {
		name           string
		msg            rawMessage
		expectedOutput []byte
		expectedError  bool
		returnError    bool
	}{
		{
			name: "successful block",
			msg: rawMessage{
				Typ: "block",
				Msg: []byte(`{"typ":"block","requserid":1,"userid":1,"blockedIDs":[]}`),
			},
			expectedOutput: []byte(`{"typ":"blocked","blockedIDs":[],"user":"nil","current":1,"tweetObjects":[]}`),
			expectedError:  false,
			returnError:    false,
		},
		{
			name: "successful unblock",
			msg: rawMessage{
				Typ: "unblock",
				Msg: []byte(`{"typ":"unblock","requserid":1,"userid":0}`),
			},
			expectedOutput: []byte(`{"typ":"unblock","userid":0,"tweetObjects":[]}`),
			expectedError:  false,
			returnError:    false,
		},
		{
			name: "failed block",
			msg: rawMessage{
				Typ: "block",
				Msg: []byte(`{"typ":"block","requserid":1,"userid":0,"blockedIDs":[]}`),
			},
			expectedOutput: []byte(`{"typ":"blocked","blockedIDs":[],"user":"nil","current":0,"tweetObjects":[]}`),
			expectedError:  false,
			returnError:    true,
		},
		{
			name: "failed unblock",
			msg: rawMessage{
				Typ: "unblock",
				Msg: []byte(`{"typ":"unblock","requserid":1,"userid":0}`),
			},
			expectedOutput: []byte(`{"typ":"unblock","userid":0,"tweetObjects":null}`),
			expectedError:  false,
			returnError:    true,
		},
	}

	m := mockDB{}
	testObj := NewBlockHandler(&m)
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

func Test_CanHandle_Block(t *testing.T) {
	m := mockDB{}
	testObj := NewBlockHandler(&m)

	testMsg := rawMessage{
		Typ: "block",
		Msg: []byte(`{"typ":"block"}`),
	}
	data := testObj.CanHandle(testMsg)
	assert.EqualValues(t, true, data, "output not as expected")

	testMsg2 := rawMessage{
		Typ: "unblock",
		Msg: []byte(`{"typ":"unblock"}`),
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
