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

type Login struct {
	Uid      int    `json:"id"`
	Typ      string `json:"typ"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
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
	Login(string, string) (int, bool, error)
	Register(string, string) (bool, error)
	CheckLike(int, int) (bool, error)
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
	`)
	return db, nil
}

func (d database) Login(username string, password string) (int, bool, error) {
	rows, _ := d.database.Query(`SELECT Id, Username, Credentials FROM UserData WHERE Username=$1`, username)
	pass := tempPass{}
	rows.Next()
	rows.Scan(&pass.Uid, &pass.username, &pass.password)

	err := bcrypt.CompareHashAndPassword([]byte(pass.password), []byte(password))
	if err != nil {
		return pass.Uid, false, nil
	} else if err == nil && username == pass.username {
		return pass.Uid, true, nil
	}

	return pass.Uid, false, nil
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
	result, err := d.database.Query(`SELECT Id, UserID FROM LikedTweets WHERE Id=$1 AND UserID=$2;`, tweetid, userid)
	if err != nil {
		return false, err
	}

	if result.Next() {
		return false, nil
	}
	return true, nil
}
