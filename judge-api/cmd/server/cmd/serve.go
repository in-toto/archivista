package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect/sql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	judgeapi "github.com/testifysec/judge/judge-api"
	"github.com/testifysec/judge/judge-api/ent"
	"github.com/testifysec/judge/judge-api/internal/auth"
	"github.com/testifysec/judge/judge-api/internal/configuration"
	"github.com/testifysec/judge/judge-api/internal/database/mysqlstore"
	policy_decision "github.com/testifysec/judge/judge-api/policy/policy_decision"
)

// This struct represents our JudgeApiServer that we create from ent.
// it is helpful in our code structure for test harnessing
type JudgeApiServer struct {
	authProvider   *auth.KratosAuthProvider
	authMiddleware mux.MiddlewareFunc
	srv            http.Handler
	database       *ent.Client
	mysqlStoreCh   <-chan error
	err            error
}

var (
	Config configuration.Config
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run:   Run,
}

func init() {
	serveCmd.PersistentFlags().StringVar(&Config.ListenOn, "listen", "tcp://127.0.0.1:8080", "Address to listen on")
	serveCmd.PersistentFlags().StringVar(&Config.SQLStoreConnectionString, "sql-connection", "", "SQL store connection string")
	serveCmd.PersistentFlags().BoolVar(&Config.GraphqlWebClientEnable, "graphql-web-client", true, "Enable the GraphQL web client")
	serveCmd.PersistentFlags().StringSliceVar(&Config.CORSAllowOrigins, "cors-origins", []string{}, "Allowed CORS origins")
	serveCmd.PersistentFlags().StringVar(&Config.KratosAdminUrl, "kratos-admin-url", "https://kratos-admin.testifysec.localhost", "Kratos admin url")
	logrus.SetFormatter(&nested.Formatter{})
}

// This func let's us set up our database independent from our judgeapiserver,
// this is crucial in testing so that we can use an in-memory db in place of a full-blown sql instance
func SetupDb(ctx context.Context, drv *sql.Driver) JudgeApiServer {

	logrus.Infof("Starting...")

	mysqlStore, mysqlStoreCh, err := mysqlstore.New(ctx, Config, drv)
	if err != nil {
		logrus.Fatalf("failed to create mysql store: %v", err)
		return JudgeApiServer{err: err}
	}

	client := mysqlStore.GetClient()
	authProvider := auth.NewKratosAuthProvider()
	authMiddleware := auth.Middleware(authProvider)

	srv := handler.NewDefaultServer(judgeapi.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	return JudgeApiServer{
		authProvider:   authProvider,
		authMiddleware: authMiddleware,
		srv:            srv,
		database:       client,
		mysqlStoreCh:   mysqlStoreCh,
	}
}

// This runs the Judge server
func Run(cmd *cobra.Command, args []string) {
	ctx, cancel := signal.NotifyContext(
		cmd.Context(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()
	startTime := time.Now()

	dbConfig, err := mysql.ParseDSN(Config.SQLStoreConnectionString)
	if err != nil {
		return
	}

	dbConfig.ParseTime = true
	Config.SQLStoreConnectionString = dbConfig.FormatDSN()
	drv, err := sql.Open("mysql", Config.SQLStoreConnectionString)
	if err != nil {
		return
	}

	judge := SetupDb(ctx, drv)
	authProvider := judge.authProvider
	authMiddleware := judge.authMiddleware
	mysqlStoreCh := judge.mysqlStoreCh
	srv := judge.srv
	client := judge.database
	router := SetupRouting(authProvider, authMiddleware, srv, client, Config)

	logrus.Infof("Serving on %s", Config.ListenOn)

	listenAddress := Config.ListenOn
	listenAddress = strings.TrimSpace(listenAddress)

	proto := ""
	if strings.HasPrefix(listenAddress, "tcp://") {
		proto = "tcp"
		listenAddress = listenAddress[len("tcp://"):]
	} else if strings.HasPrefix(listenAddress, "unix://") {
		proto = "unix"
		listenAddress = listenAddress[len("unix://"):]
	}

	listener, err := net.Listen(proto, listenAddress)
	if err != nil {
		logrus.Fatalf("unable to start http listener: %+v", err)
	}

	server := &http.Server{
		Addr: listenAddress,
		Handler: handlers.CORS(
			handlers.AllowedOrigins(Config.CORSAllowOrigins),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"}),
		)(router),
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("unable to start http server: %+v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()
		if err := server.Shutdown(ctxShutdown); err != nil {
			logrus.Errorf("server shutdown failed: %+v", err)
		}
	}()

	logrus.Infof("startup complete (time since start: %s)", time.Since(startTime))
	<-ctx.Done()
	<-mysqlStoreCh
}

// This sets up the Routing for the Judge Server
func SetupRouting(authProvider *auth.KratosAuthProvider, authMiddleware mux.MiddlewareFunc, srv http.Handler, database *ent.Client, config configuration.Config) *mux.Router {
	router := mux.NewRouter()

	authSubrouter := router.PathPrefix("/").Subrouter()
	authSubrouter.Use(authMiddleware)

	authSubrouter.Handle("/query", srv)
	if config.GraphqlWebClientEnable {
		authSubrouter.Handle("/",
			playground.Handler("Judge", "/query"),
		)
	}

	// WebhookSubrouter does not have cookie auth middleware
	webhookSubrouter := router.PathPrefix("/webhook").Subrouter()
	webhookSubrouter.Handle("/defaulttenant", http.HandlerFunc(authProvider.UpdateAssignedTenantsWithIdentityId)).Methods(http.MethodPost)

	router.HandleFunc("/policy-decision", func(w http.ResponseWriter, r *http.Request) {
		policy_decision.PostPolicy(w, r, database)
	})

	return router
}
