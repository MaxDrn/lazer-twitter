package persistence

import (
	"database/sql"
	"fmt"
	options "lazer-twitter/pkg/cli/lazer-twitter"

	"github.com/fid-dev/go-pflog/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type ClientTweet struct {
	Id      int    `json:"id"`
	Time    string `json:"time"`
	Likes   int    `json:"likes"`
	UserID  int    `json:"userid"`
	User    string `json:"user"`
	Message string `json:"message"`
}

type likedMessage struct {
	Typ     string `json:"typ"`
	UserID  int    `json:"userid"`
	TweetId int    `json:"tweetid"`
}

type Login struct {
	Uid              int           `json:"id"`
	Typ              string        `json:"typ"`
	Username         string        `json:"username"`
	BlockedIds       []int         `json:"blockedids"`
	BlockedUsernames []string      `json:"blockedusernames"`
	Password         string        `json:"password,omitempty"`
	UpdatedTweets    []ClientTweet `json:"tweetObjects"`
}

func NewDatabase() (*database, error) {
	db, err := connectToDatabase()
	return &database{
		database: db,
	}, err
}

type Database interface {
	InsertTweet(*ClientTweet) (int, error)
	GetAllTweets() ([]ClientTweet, error)
	LikeTweet(int, int) error
	GetTweet(int) (*ClientTweet, error)
	Login(string, string) (int, []int, bool, error)
	Register(string, string) (bool, error)
	CheckLike(int, int) (bool, error)
	FilteredTweets(int) ([]ClientTweet, string, error)
	GetTweetsFromUserID(int) ([]ClientTweet, error)
	InsertBlockedUser(int, int) (bool, error)
	RemoveBlockedUser(int, int) (bool, error)
	GetBlockedUserFromId(int) ([]int, error)
}
type database struct {
	database *sql.DB
}

type tempPass struct {
	Uid      int
	username string
	password string
}

var _ Database = &database{}

func (d database) InsertTweet(inf *ClientTweet) (int, error) {
	lastInsertId := 0
	err := d.database.QueryRow(`INSERT INTO Tweets VALUES(DEFAULT, $1, $2, $3, $4, $5) RETURNING id`, inf.Time, inf.Likes, inf.UserID, inf.User, inf.Message).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (d database) GetAllTweets() ([]ClientTweet, error) {
	result, err := d.database.Query(`
	SELECT * FROM Tweets ORDER BY Id;
	`)
	tweets := make([]ClientTweet, 0)
	for result.Next() {
		copyTweet := ClientTweet{}
		err := result.Scan(&copyTweet.Id, &copyTweet.Time, &copyTweet.Likes, &copyTweet.UserID, &copyTweet.User, &copyTweet.Message)
		if err != nil {
			log.Errorf("could not find tweets %v", err)
			continue
		}
		tweets = append(tweets, copyTweet)
	}

	if err != nil {
		return nil, err
	}
	return tweets, nil
}

func (d database) LikeTweet(id int, userID int) error {
	_, err := d.database.Exec(`
	UPDATE Tweets SET Likes=Likes+1 WHERE Id=$1;
	`, id)
	if err != nil {
		return errors.Wrap(err, "failed to up the like count")
	}

	_, err = d.database.Exec(`
	INSERT INTO LikedTweets VALUES($1, $2);
	`, id, userID)
	if err != nil {
		return errors.Wrap(err, "failed to up the like count")
	}

	return nil
}

func (d database) GetTweet(id int) (*ClientTweet, error) {
	result, err := d.database.Query(`
		SELECT * FROM Tweets WHERE Id=$1;
	`, id)

	if err != nil {
		return nil, err
	}
	inf := ClientTweet{}

	if result.Next() {
		err := result.Scan(&inf.Id, &inf.Time, &inf.Likes, &inf.UserID, &inf.User, &inf.Message)

		if err != nil {
			return nil, err
		}
		return &inf, nil
	}

	return nil, nil
}

func connectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable",
		options.Current.DBUser,
		options.Current.DBName,
		options.Current.DBPassword,
		options.Current.DBPort))

	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to database %s", options.Current.DBName)
	}

	err = db.Ping()

	if err != nil {
		return nil, errors.Wrapf(err, "could not ping database %s", options.Current.DBName)
	} else {
		log.Info("successfully connected to the database")
	}

	db.Exec(`
	CREATE TABLE IF NOT EXISTS Tweets(Id SERIAL PRIMARY KEY, TweetTime text, Likes int, UserID int, Username text, Message text);
	CREATE TABLE IF NOT EXISTS UserData(Id SERIAL PRIMARY KEY, Username text, Credentials text);
	CREATE TABLE IF NOT EXISTS LikedTweets(Id int, UserID int);
	CREATE TABLE IF NOT EXISTS BlockedUser(Id int, BlockedID int);
	`)
	return db, nil
}

func (d database) Login(username string, password string) (int, []int, bool, error) {
	rows, _ := d.database.Query(`SELECT Id, Username, Credentials FROM UserData WHERE Username=$1`, username)
	pass := tempPass{}
	rows.Next()
	_ = rows.Scan(&pass.Uid, &pass.username, &pass.password)

	err := bcrypt.CompareHashAndPassword([]byte(pass.password), []byte(password))
	if err != nil {
		return pass.Uid, nil, false, nil
	} else if err == nil && username == pass.username {
		blockedUser, err := d.GetBlockedUserFromId(pass.Uid)
		if err != nil {
			return pass.Uid, nil, false, err
		}
		return pass.Uid, blockedUser, true, nil
	}

	return pass.Uid, nil, false, nil
}

func (d database) Register(username string, password string) (bool, error) {
	result, err := d.database.Query(`SELECT Username from UserData WHERE Username=$1;`, username)

	if result.Next() {
		return false, errors.Wrap(err, "could not register into database")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return false, errors.Wrap(err, "could not hash password")
	}
	_, err = d.database.Exec(`INSERT INTO UserData Values(DEFAULT, $1, $2);`, username, hash)
	if err != nil {
		return false, errors.Wrap(err, "could not register into database")
	}
	return true, nil
}

func (d database) CheckLike(tweetid int, userid int) (bool, error) {
	rows, err := d.database.Query(`SELECT UserID FROM Tweets WHERE Id=$1;`, tweetid)
	_ = rows.Next()
	tempLike := likedMessage{}
	err = rows.Scan(&tempLike.UserID)
	if err != nil {
		return false, err
	}
	result, errTwo := d.database.Query(`SELECT Id, UserID FROM LikedTweets WHERE Id=$1 AND UserID=$2;`, tweetid, userid)
	if errTwo != nil {
		return false, err
	}
	hasResult := result.Next()
	if hasResult {
		return false, nil
	}
	if !hasResult && userid != tempLike.UserID {
		return true, nil
	}
	return false, nil
}

func (d database) FilteredTweets(blockid int) ([]ClientTweet, string, error) {
	tweets := make([]ClientTweet, 0)
	blockedTweet := ClientTweet{}
	userRow, err := d.database.Query(`SELECT Username FROM Tweets WHERE UserID=$1;`, blockid)
	if err != nil {
		return nil, "nil", err
	}
	userRow.Next()
	errT := userRow.Scan(&blockedTweet.User)
	if errT != nil {
		return nil, "nil", errT
	}
	result, err := d.database.Query(`SELECT * FROM Tweets WHERE NOT UserID=$1;`, blockid)
	if err != nil {
		return nil, blockedTweet.User, err
	}
	for result.Next() {
		tweet := ClientTweet{}
		err := result.Scan(&tweet.Id, &tweet.Time, &tweet.Likes, &tweet.UserID, &tweet.User, &tweet.Message)
		if err != nil {
			return nil, blockedTweet.User, err
		}
		tweets = append(tweets, tweet)
	}
	return tweets, blockedTweet.User, nil
}

func (d database) GetTweetsFromUserID(userid int) ([]ClientTweet, error) {
	tweets := make([]ClientTweet, 0)
	result, err := d.database.Query(`SELECT * FROM Tweets WHERE UserID=$1;`, userid)
	if err != nil {
		return nil, err
	}
	for result.Next() {
		tweet := ClientTweet{}
		err := result.Scan(&tweet.Id, &tweet.Time, &tweet.Likes, &tweet.UserID, &tweet.User, &tweet.Message)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

func (d database) InsertBlockedUser(userID int, blockedUser int) (bool, error) {
	_, err := d.database.Exec(`INSERT INTO BlockedUser VALUES($1, $2);`, userID, blockedUser)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (d database) RemoveBlockedUser(userID int, blockedUser int) (bool, error) {
	_, err := d.database.Exec(`DELETE FROM BlockedUser WHERE Id=$1 AND BlockedID=$2;`, userID, blockedUser)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (d database) GetBlockedUserFromId(id int) ([]int, error) {
	result, err := d.database.Query(`SELECT BlockedID FROM BlockedUser WHERE Id=$1;`, id)
	if err != nil {
		return nil, err
	}
	blockedIds := make([]int, 0)
	temp := Login{}
	for result.Next() {
		err := result.Scan(&temp.Uid)
		if err != nil {
			return nil, err
		}
		blockedIds = append(blockedIds, temp.Uid)
	}
	return blockedIds, nil
}
