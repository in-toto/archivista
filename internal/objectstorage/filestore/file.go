package filestore

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/git-bom/gitbom-go"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/archivist-api/pkg/api/archivist"
	"google.golang.org/protobuf/types/known/emptypb"
)

type store struct {
	archivist.UnimplementedCollectorServer

	prefix string
}

func NewServer(ctx context.Context, directory string, address string) (archivist.CollectorServer, <-chan error, error) {
	errCh := make(chan error)
	go func() {

		server := handlers.CompressHandler(http.FileServer(http.Dir(directory)))
		log.Fatalln(http.ListenAndServe(address, server))

		<-ctx.Done()
		close(errCh)
	}()

	return &store{
		prefix: directory,
	}, errCh, nil
}

func (s *store) Get(ctx context.Context, request *archivist.GetRequest) (*archivist.GetResponse, error) {
	res, err := ioutil.ReadFile(filepath.Join(s.prefix, request.GetGitoid()))
	if err != nil {
		logrus.WithContext(ctx).Errorf("failed to retrieve object: %+v", err)
		return nil, err
	}
	return &archivist.GetResponse{
		Gitoid: request.GetGitoid(),
		Object: string(res),
	}, nil
}

func (s *store) Store(ctx context.Context, request *archivist.StoreRequest) (*emptypb.Empty, error) {
	// TODO refactor this to use common code
	gb := gitbom.NewSha256GitBom()
	if err := gb.AddReference([]byte(request.Object), nil); err != nil {
		logrus.WithContext(ctx).Errorf("gitbom tag generation failed: %+v", err)
		return nil, err
	}

	fmt.Printf("Writing: %s/%s\n", s.prefix, gb.Identity())
	err := ioutil.WriteFile(filepath.Join(s.prefix, gb.Identity()+".json"), []byte(request.Object), 0644)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
