package http

import (
	"lazer-twitter/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ persistence.Database = &mockDB{}

func Test_MessageHandler(t *testing.T) {

	cases := []struct {
		name           string
		msg            rawMessage
		expectedOutput []byte
		expectedError  bool
	}{
		{
			name: "Default test",
			msg: rawMessage{
				Typ: "message",
				Msg: nil,
			},
			expectedOutput: []byte(`{"typ":"","tweet":{"id":0,"time":"","likes":0,"user":"","message":""}}`),
			expectedError:  false,
		},
	}

	m := mockDB{}
	testObj := NewMessageHandler(&m)
	for _, val := range cases {
		t.Run(val.name, func(tt *testing.T) {
			tweets, _, err := testObj.Handle(val.msg)
			assert.EqualValues(tt, val.expectedOutput, tweets, "Output not as expected")
			assert.EqualValues(tt, val.expectedError, err != nil, "Error not as expected")
		})
	}
}

func Test_CanHandle_Message(t *testing.T) {
	m := mockDB{}
	testObj := NewMessageHandler(&m)

	testMsg := rawMessage{
		Typ: "message",
		Msg: []byte(`{"typ":"message"}`),
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
