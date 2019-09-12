package http

import (
	"lazer-twitter/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ persistence.Database = &mockDB{}

func Test_LikeHandler(t *testing.T) {

	cases := []struct {
		name              string
		msg               rawMessage
		expectedOutput    []byte
		expectedError     bool
		expectedLikeCalls int
		returnError       bool
	}{
		{
			name: "successful like",
			msg: rawMessage{
				Typ: "like",
				Msg: []byte(`{"typ":"like","userid":0,"tweetid":0}`),
			},
			expectedOutput:    []byte(`{"typ":"liked","tweet":{"id":0,"time":"","likes":0,"userid":0,"user":"","message":""}}`),
			expectedError:     false,
			expectedLikeCalls: 1,
			returnError:       false,
		},
		{
			name: "failed like",
			msg: rawMessage{
				Typ: "like",
				Msg: []byte(`{"typ":"like","userid":0,"tweetid":"0"`),
			},
			expectedOutput:    []byte(`{"typ":"failedLike","userid":0,"user":"","tweetid":0}`),
			expectedError:     false,
			expectedLikeCalls: 1,
			returnError:       true,
		},
	}

	m := mockDB{}
	testObj := NewLikeHandler(&m)
	for _, val := range cases {
		t.Run(val.name, func(tt *testing.T) {
			if val.returnError == true {
				m.returnError = true
			}
			tweets, _, err := testObj.Handle(val.msg)
			assert.EqualValues(tt, val.expectedLikeCalls, m.likedCalls, "call count not as expected")
			assert.EqualValues(tt, val.expectedOutput, tweets, "Output not as expected")
			assert.EqualValues(tt, val.expectedError, err != nil, "Error not as expected")
		})
		m.likedCalls = 0
		m.returnError = false
	}
}

func Test_CanHandle_Like(t *testing.T) {
	m := mockDB{}
	testObj := NewLikeHandler(&m)

	testMsg := rawMessage{
		Typ: "like",
		Msg: []byte(`{"typ":"like"}`),
	}
	data := testObj.CanHandle(testMsg)
	assert.EqualValues(t, true, data, "output not as expected")

	testMsg2 := rawMessage{
		Typ: "join",
		Msg: []byte(`{"typ":"join"}`),
	}

	dataTwo := testObj.CanHandle(testMsg2)
	assert.EqualValues(t, false, dataTwo, "output not as expected")
}
