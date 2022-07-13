package filestore

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
)

type Store struct {
	prefix string
}

func New(ctx context.Context, directory string, address string) (*Store, <-chan error, error) {
	errCh := make(chan error)
	go func() {
		server := handlers.CompressHandler(http.FileServer(http.Dir(directory)))
		log.Fatalln(http.ListenAndServe(address, server))
		<-ctx.Done()
		close(errCh)
	}()

	return &Store{
		prefix: directory,
	}, errCh, nil
}

func (s *Store) Get(ctx context.Context, request *archivist.GetRequest) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.prefix, request.GetGitoid()+".json"))
}

func (s *Store) Store(ctx context.Context, gitoid string, payload []byte) error {
	return os.WriteFile(filepath.Join(s.prefix, gitoid+".json"), payload, 0644)
}
