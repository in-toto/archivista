package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	archivistGrpcUrl string
	archivistGqlUrl  string

	rootCmd = &cobra.Command{
		Use:   "archivistctl",
		Short: "A utility to interact with an archivist server",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&archivistGrpcUrl, "grpc-url", "u", "localhost:8080", "url of the archivist grpc service")
	rootCmd.PersistentFlags().StringVarP(&archivistGqlUrl, "graphql-url", "q", "http://localhost:8082/query", "url of the archivist graphql service")
}

func Execute() error {
	return rootCmd.Execute()
}

func newConn(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
