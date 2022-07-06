package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	archivistUrl string

	rootCmd = &cobra.Command{
		Use:   "archivistctl",
		Short: "A utility to interact with an archivist server",
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&archivistUrl, "archivisturl", "u", "localhost:8080", "url of the archivist instance")
}

func Execute() error {
	return rootCmd.Execute()
}

func newConn(url string) (*grpc.ClientConn, error) {
	return grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
