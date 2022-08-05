package archivist

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/testifysec/archivist/ent"
)

func (r *queryResolver) Status(ctx context.Context) (*Status, error) {
	_, err := r.client.Dsse.Query().First(ctx)
	databaseStatus := "offline"
	if err == nil {
		databaseStatus = "online"
	}
	return &Status{
		Database: databaseStatus,
		Graphql:  "online",
		Grpc:     "online",
	}, nil
}

func (r *queryResolver) Dsses(ctx context.Context, after *ent.Cursor, first *int, before *ent.Cursor, last *int, where *ent.DsseWhereInput) (*ent.DsseConnection, error) {
	return r.client.Dsse.Query().
		Paginate(ctx, after, first, before, last,
			ent.WithDsseFilter(where.Filter))
}

func (r *statementResolver) AttestationCollection(ctx context.Context, obj *ent.Statement) (*ent.AttestationCollection, error) {
	return obj.AttestationCollections(ctx)
}

func (r *subjectResolver) SubjectDigest(ctx context.Context, obj *ent.Subject) ([]*ent.SubjectDigest, error) {
	return obj.SubjectDigests(ctx)
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Statement returns StatementResolver implementation.
func (r *Resolver) Statement() StatementResolver { return &statementResolver{r} }

// Subject returns SubjectResolver implementation.
func (r *Resolver) Subject() SubjectResolver { return &subjectResolver{r} }

type queryResolver struct{ *Resolver }
type statementResolver struct{ *Resolver }
type subjectResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
type dsseResolver struct{ *Resolver }
