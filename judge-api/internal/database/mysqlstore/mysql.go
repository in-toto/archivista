package mysqlstore

import (
	"context"
	"fmt"
	"time"

	"ariga.io/sqlcomment"
	"entgo.io/ent/dialect/sql"
	"github.com/sirupsen/logrus"
	"github.com/testifysec/judge/judge-api/ent"
	"github.com/testifysec/judge/judge-api/internal/configuration"
)

const subjectBatchSize = 30000

// mysql has a limit of 65536 parameters in a single query. each subject has ~3 parameters [subject id, algo, value],
// so we can theoretically jam 65535/3 subjects in a single batch. but we probably want some breathing room just in case.
const subjectDigestBatchSize = 20000

type Store struct {
	client *ent.Client
}

func New(ctx context.Context, cfg configuration.Config, drv *sql.Driver) (*Store, <-chan error, error) {
	sqlcommentDrv := sqlcomment.NewDriver(drv,
		sqlcomment.WithDriverVerTag(),
		sqlcomment.WithTags(sqlcomment.Tags{
			sqlcomment.KeyApplication: "judge-api",
			sqlcomment.KeyFramework:   "net/http",
		}),
	)

	db := drv.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(3 * time.Minute)

	client := ent.NewClient(ent.Driver(sqlcommentDrv), ent.Config(cfg))

	errCh := make(chan error)

	go func() {
		<-ctx.Done()
		err := client.Close()
		if err != nil {
			logrus.Errorf("error closing database: %+v", err)
		}
		close(errCh)
	}()

	if err := client.Schema.Create(ctx); err != nil {
		logrus.Fatalf("failed creating schema resources: %v", err)
	}

	return &Store{
		client: client,
	}, errCh, nil
}

func (s *Store) withTx(ctx context.Context, fn func(tx *ent.Tx) error) error {
	tx, err := s.client.Tx(ctx)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("unable to rollback transaction: %w", err)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}

func (s *Store) Store(ctx context.Context, obj []byte) error {
	return s.withTx(ctx, func(tx *ent.Tx) error {
		_, err := tx.Project.Create().Save(ctx)
		return err
	})
}

func (s *Store) GetClient() *ent.Client {
	return s.client
}

type saver[T any] interface {
	Save(context.Context) ([]T, error)
}

func batch[TCreate any, TResult any](ctx context.Context, batchSize int, create []TCreate, saveFn func(...TCreate) saver[TResult]) ([]TResult, error) {
	results := make([]TResult, 0, len(create))
	for i := 0; i < len(create); i += batchSize {
		var batch []TCreate
		if i+batchSize > len(create) {
			batch = create[i:]
		} else {
			batch = create[i : i+batchSize]
		}

		batchSaver := saveFn(batch...)
		batchResults, err := batchSaver.Save(ctx)
		if err != nil {
			return nil, err
		}

		results = append(results, batchResults...)
	}

	return results, nil
}
