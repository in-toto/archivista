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

	"gitlab.com/testifysec/judge-platform/judge-api/internal/auth"
	"gitlab.com/testifysec/judge-platform/judge-api/internal/configuration"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/networkservicemesh/sdk/pkg/tools/debug"
	"github.com/networkservicemesh/sdk/pkg/tools/log"
	"github.com/networkservicemesh/sdk/pkg/tools/log/logruslogger"

	judgeapi "gitlab.com/testifysec/judge-platform/judge-api"
	"gitlab.com/testifysec/judge-platform/judge-api/internal/database/mysqlstore"
)

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
}

func Run(cmd *cobra.Command, args []string) {
	startTime := time.Now()
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	defer cancel()

	logrus.SetFormatter(&nested.Formatter{})
	log.EnableTracing(true)
	ctx = log.WithLog(ctx, logruslogger.New(ctx))

	if err := debug.Self(); err != nil {
		log.FromContext(ctx).Infof("%s", err)
	}

	log.FromContext(ctx).Infof("Starting...")
	mysqlStore, mysqlStoreCh, err := mysqlstore.New(ctx, Config)
	if err != nil {
		log.FromContext(ctx).Fatalf("failed to create mysql store: %v", err)
	}

	client := mysqlStore.GetClient()
	authProvider := auth.NewKratosAuthProvider()
	authMiddleware := auth.Middleware(authProvider)

	srv := handler.NewDefaultServer(judgeapi.NewSchema(client))
	srv.Use(entgql.Transactioner{TxOpener: client})

	router := mux.NewRouter()
	router.Use(authMiddleware)
	router.Handle("/query", srv)
	if Config.GraphqlWebClientEnable {
		router.Handle("/",
			playground.Handler("Judge", "/query"),
		)
	}

	log.FromContext(ctx).Infof("Serving on %s", Config.ListenOn)

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
		log.FromContext(ctx).Fatalf("unable to start http listener: %+v", err)
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
			log.FromContext(ctx).Fatalf("unable to start http server: %+v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()
		if err := server.Shutdown(ctxShutdown); err != nil {
			log.FromContext(ctx).Errorf("server shutdown failed: %+v", err)
		}
	}()

	log.FromContext(ctx).Infof("startup complete (time since start: %s)", time.Since(startTime))

	<-ctx.Done()
	<-mysqlStoreCh

	log.FromContext(ctx).Infof("exiting, uptime: %v", time.Since(startTime))

}
