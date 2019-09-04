package cmd

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	options "lazer-twitter/pkg/cli/lazer-twitter"
	"lazer-twitter/pkg/http"
	netHttp "net/http"
)

func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "lazer-twitter",
		Short: "lazer-twitter is a simple Twitter Clone.",
		Run: func(cmd *cobra.Command, args []string) {

			db := connectToDatabase()

			http.Database = db

			netHttp.Handle("/socket", netHttp.HandlerFunc(http.SocketHandler))
			netHttp.Handle("/", netHttp.FileServer(netHttp.Dir("./assets")))
			_ = netHttp.ListenAndServe("localhost:"+options.Current.RESTListenPort, nil)
		},
	}

	cmd.Flags().StringVar(&options.Current.RESTListenPort, "rest-listen-port", options.Current.RESTListenPort, "tcp port to listen for HTTP requests")
	cmd.Flags().StringVar(&options.Current.DBName, "db-name", options.Current.DBName, "database name")
	cmd.Flags().StringVar(&options.Current.DBUser, "db-user", options.Current.DBUser, "database user")
	cmd.Flags().StringVar(&options.Current.DBPassword, "db-pw", options.Current.DBPassword, "database pw")
	cmd.Flags().StringVar(&options.Current.DBHost, "db-host", options.Current.DBHost, "database host")
	cmd.Flags().StringVar(&options.Current.DBPort, "db-port", options.Current.DBPort, "database port")
	return cmd
}

func connectToDatabase() *sql.DB{
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s port=%s sslmode=disable",
		options.Current.DBUser,
		options.Current.DBName,
		options.Current.DBPassword,
		options.Current.DBPort))

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()

	if err != nil {
		panic(err.Error())
		fmt.Println("could not connect to the database")
		return nil
	} else {
		fmt.Println("successfully connected to the database!")
	}

	db.Exec(`
	CREATE TABLE IF NOT EXISTS Tweets(Typ text, TweetTime text, Likes int, UserName text, Message text, TweetID int);
	`)
	return db
}
