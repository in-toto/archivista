package publisher

import (
	"context"

	"github.com/in-toto/archivista/pkg/message"
)

type Publisher interface {
	Publish(ctx context.Context, msg message.Message) error
}
