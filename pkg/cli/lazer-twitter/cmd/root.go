package cmd

import (
	"lazer-twitter/persistence"
	options "lazer-twitter/pkg/cli/lazer-twitter"
	"lazer-twitter/pkg/http"
	netHttp "net/http"

	"github.com/fid-dev/go-pflog/log"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "lazer-twitter",
		Short: "lazer-twitter is a simple Twitter Clone.",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := persistence.NewDatabase()

			if err != nil {
				log.Fatal(err.Error())
			}

			handler := http.NewWebSocketHandler(db)

			netHttp.Handle("/socket", netHttp.HandlerFunc(handler.SocketHandler))
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
