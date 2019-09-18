package http

import (
	"fmt"
	"lazer-twitter/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	returnError     bool
	likedCalls      int
	returnedObjects int
}

func (m *mockDB) InsertBlockedUser(int, int) (bool, error) {
	if m.returnError {
		return false, nil
	} else if !m.returnError {
		return true, nil
	}
	return true, nil
}

func (m *mockDB) RemoveBlockedUser(int, int) (bool, error) {
	if m.returnError {
		return false, nil
	} else if !m.returnError {
		return true, nil
	}
	return true, nil
}

func (m *mockDB) GetBlockedUserFromId(int) ([]int, error) {
	return nil, nil
}

func (m *mockDB) FilteredTweets(int) ([]persistence.ClientTweet, string, error) {
	fmt.Println("test")
	return nil, "nil", nil
}

func (m *mockDB) GetTweetsFromUserID(int) ([]persistence.ClientTweet, error) {
	if m.returnError {
		return nil, nil
	} else if !m.returnError {
		return []persistence.ClientTweet{}, nil
	}
	return []persistence.ClientTweet{}, nil
}

func (m *mockDB) InsertTweet(tweet *persistence.ClientTweet) (int, error) {
	return 0, nil
}

func (m *mockDB) GetTweet(int) (*persistence.ClientTweet, error) {
	return nil, nil
}

func (m *mockDB) Login(string, string) (int, []int, bool, error) {
	if m.returnError == true {
		return 0, nil, false, nil
	} else {
		return 0, nil, true, nil
	}
}

func (m *mockDB) Register(string, string) (bool, error) {
	if m.returnError == true {
		return false, nil
	} else {
		return true, nil
	}
}

func (m *mockDB) CheckLike(int, int) (bool, error) {
	m.likedCalls++
	if m.returnError == true {
		return false, nil
	} else {
		return true, nil
	}
}

var _ persistence.Database = &mockDB{}

func (m *mockDB) InsertIntoDatabase(tweet *persistence.ClientTweet) (int, error) {
	return 0, nil
}
func (m *mockDB) GetAllTweets() ([]persistence.ClientTweet, error) {
	mockTweet := persistence.ClientTweet{
		Id:      0,
		Time:    "now",
		Likes:   2,
		UserID:  0,
		User:    "Max",
		Message: "Hallo",
	}

	mockTweet2 := persistence.ClientTweet{
		Id:      1,
		Time:    "now",
		Likes:   1,
		UserID:  0,
		User:    "Peter",
		Message: "Hey",
	}
	if m.returnedObjects == 0 {
		return nil, nil
	} else if m.returnedObjects == 1 {
		return []persistence.ClientTweet{mockTweet}, nil
	} else if m.returnedObjects == 2 {
		return []persistence.ClientTweet{mockTweet, mockTweet2}, nil
	}
	return nil, nil
}
func (m *mockDB) LikeTweet(i int, ii int) error {
	return nil
}
func (m *mockDB) GetRow(j int) (*persistence.ClientTweet, error) {
	mockTweet := persistence.ClientTweet{
		Id:      1,
		Time:    "now",
		Likes:   1,
		UserID:  0,
		User:    "Peter",
		Message: "Hey",
	}
	return &mockTweet, nil
}

func Test_JoinHandler(t *testing.T) {

	cases := []struct {
		name            string
		msg             rawMessage
		returnedObjects int
		expectedOutput  []byte
		expectedError   bool
	}{
		{
			name: "Empty reply",
			msg: rawMessage{
				Typ: "join",
				Msg: []byte(`{"typ":"join"}`),
			},
			returnedObjects: 0,
			expectedOutput:  []byte(`{"typ":"all","tweetObjects":null}`),
			expectedError:   false,
		},
		{
			name: "Single reply",
			msg: rawMessage{
				Typ: "join",
				Msg: []byte(`{"typ":"join"}`),
			},
			returnedObjects: 1,
			expectedOutput:  []byte(`{"typ":"all","tweetObjects":[{"id":0,"time":"now","likes":2,"userid":0,"user":"Max","message":"Hallo"}]}`),
			expectedError:   false,
		},
		{
			name: "multiple reply - 2",
			msg: rawMessage{
				Typ: "join",
				Msg: []byte(`{"typ":"join"}`),
			},
			returnedObjects: 2,
			expectedOutput:  []byte(`{"typ":"all","tweetObjects":[{"id":0,"time":"now","likes":2,"userid":0,"user":"Max","message":"Hallo"},{"id":1,"time":"now","likes":1,"userid":0,"user":"Peter","message":"Hey"}]}`),
			expectedError:   false,
		},
	}

	m := mockDB{}
	testObj := NewJoinHandler(&m)
	for _, val := range cases {
		t.Run(val.name, func(tt *testing.T) {
			m.returnedObjects = val.returnedObjects
			tweets, _, err := testObj.Handle(val.msg)
			assert.EqualValues(tt, val.expectedOutput, tweets, "Output not as expected")
			assert.EqualValues(tt, val.expectedError, err != nil, "Error not as expected")
		})
	}
}

func Test_CanHandle_Join(t *testing.T) {
	m := mockDB{}
	testObj := NewJoinHandler(&m)

	testMsg := rawMessage{
		Typ: "join",
		Msg: []byte(`{"typ":"join"}`),
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
