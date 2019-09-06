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

func NewDatabase() (*database, error) {
	db, err := connectToDatabase()
	return &database{
		database: db,
	}, err
}

type Database interface {
	InsertIntoDatabase(tweet *ClientTweet) (int, error)
	GetAllTweets() ([]ClientTweet, error)
	LikeTweet(int)
	GetRow(int) (*ClientTweet, error)
}
type database struct {
	database *sql.DB
}

var _ Database = &database{}

func (d database) InsertIntoDatabase(inf *ClientTweet) (int, error) {
	lastInsertId := 0
	d.database.QueryRow(`INSERT INTO Tweets VALUES(DEFAULT, $1, $2, $3, $4) RETURNING id`, inf.Time, inf.Likes, inf.User, inf.Message).Scan(&lastInsertId)
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

func (d database) LikeTweet(id int) {
	_, err := d.database.Exec(`
	UPDATE Tweets SET Likes=Likes+1 WHERE Id=$1;
	`, id)

	if err != nil {
		log.Error(err.Error())
	}
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
	`)
	return db, nil
}
