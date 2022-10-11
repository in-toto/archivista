package archivist

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/testifysec/archivist/ent"
)

func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.client.Noder(ctx, id)
}

func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids)
}

func (r *queryResolver) Dsses(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, where *ent.DsseWhereInput) (*ent.DsseConnection, error) {
	return r.client.Dsse.Query().Paginate(ctx, after, first, before, last, ent.WithDsseFilter(where.Filter))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
