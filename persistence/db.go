package persistence

import (
	"database/sql"
	"fmt"
	options "lazer-twitter/pkg/cli/lazer-twitter"

	"github.com/fid-dev/go-pflog/log"
	"github.com/pkg/errors"
)

type ClientTweet struct {
	Id      int    `json:"id"`
	Time    string `json:"time"`
	Likes   int    `json:"likes"`
	User    string `json:"user"`
	Message string `json:"message"`
}

type Login struct {
	Typ      string `json:"typ"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewDatabase() (*database, error) {
	db, err := connectToDatabase()
	return &database{
		database: db,
	}, err
}

type Database interface {
	InsertIntoDatabase(tweet *ClientTweet) (int, error)
	GetAllTweets() ([]ClientTweet, error)
	LikeTweet(int, string) error
	GetRow(int) (*ClientTweet, error)
	LoginDatabase(string, string) (bool, error)
	RegisterDatabase(string, string) (bool, error)
	CheckLike(int, string) (*sql.Rows, error)
}
type database struct {
	database *sql.DB
}

var _ Database = &database{}

func (d database) InsertIntoDatabase(inf *ClientTweet) (int, error) {
	lastInsertId := 0
	err := d.database.QueryRow(`INSERT INTO Tweets VALUES(DEFAULT, $1, $2, $3, $4) RETURNING id`, inf.Time, inf.Likes, inf.User, inf.Message).Scan(&lastInsertId)
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
		err := result.Scan(&copyTweet.Id, &copyTweet.Time, &copyTweet.Likes, &copyTweet.User, &copyTweet.Message)
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

func (d database) LikeTweet(id int, username string) error {
	_, err := d.database.Exec(`
	UPDATE Tweets SET Likes=Likes+1 WHERE Id=$1;
	`, id)
	if err != nil {
		return errors.Wrap(err, "failed to up the like count")
	}

	_, err = d.database.Exec(`
	INSERT INTO LikedTweets VALUES($1, $2);
	`, id, username)
	if err != nil {
		return errors.Wrap(err, "failed to up the like count")
	}

	return nil
}

func (d database) GetRow(id int) (*ClientTweet, error) {
	result, err := d.database.Query(`
		SELECT * FROM Tweets WHERE Id=$1;
	`, id)

	if err != nil {
		return nil, err
	}
	inf := ClientTweet{}

	if result.Next() {
		err := result.Scan(&inf.Id, &inf.Time, &inf.Likes, &inf.User, &inf.Message)

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
	CREATE TABLE IF NOT EXISTS Tweets(Id SERIAL PRIMARY KEY, TweetTime text, Likes int, UserName text, Message text);
	CREATE TABLE IF NOT EXISTS UserData(Username text, Credentials text);
	CREATE TABLE IF NOT EXISTS LikedTweets(Id int, Username text);
	`)
	return db, nil
}

func (d database) LoginDatabase(username string, password string) (bool, error) {
	result, err := d.database.Query(`SELECT UserData.Username, UserData.Credentials from UserData WHERE Username=$1 AND Credentials=$2;`, username, password)
	if err != nil {
		log.Error(err.Error())
		return false, err
	}
	mockLogin := Login{}

	for result.Next() {
		err := result.Scan(&mockLogin.Username, &mockLogin.Password)
		if err != nil {
			log.Error(err.Error())
			return false, err
		}
	}

	if mockLogin.Username == username && mockLogin.Password == password {
		return true, nil
	} else {
		return false, nil
	}
}

func (d database) RegisterDatabase(username string, password string) (bool, error) {
	result, err := d.database.Query(`SELECT Username from UserData WHERE Username=$1;`, username)
	mockLogin := Login{}

	for result.Next() {
		_ = result.Scan(&mockLogin.Username)
		if mockLogin.Username == username {
			return false, errors.Wrap(err, "could not register into database")
		}
	}

	_, err = d.database.Exec(`INSERT INTO UserData Values($1, $2);`, username, password)
	if err != nil {
		return false, errors.Wrap(err, "could not register into database")
	}
	return true, nil
}

func (d database) CheckLike(id int, username string) (*sql.Rows, error) {
	result, err := d.database.Query(`SELECT Id, Username FROM LikedTweets WHERE Id=$1 AND Username=$2;`, id, username)
	if err != nil {
		return nil, errors.Wrap(err, "failed to evaluate query")
	}
	return result, nil
}
