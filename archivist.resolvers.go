// Copyright 2022 The Archivist Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
