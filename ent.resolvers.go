// Copyright 2023 The Archivista Contributors
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

package archivista

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"entgo.io/contrib/entgql"
	"github.com/testifysec/archivista/ent"
)

func (r *queryResolver) Node(ctx context.Context, id int) (ent.Noder, error) {
	return r.client.Noder(ctx, id)
}

func (r *queryResolver) Nodes(ctx context.Context, ids []int) ([]ent.Noder, error) {
	return r.client.Noders(ctx, ids)
}

func (r *queryResolver) Dsses(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, where *ent.DsseWhereInput) (*ent.DsseConnection, error) {
	return r.client.Dsse.Query().Paginate(ctx, after, first, before, last, ent.WithDsseFilter(where.Filter))
}

func (r *queryResolver) Subjects(ctx context.Context, after *entgql.Cursor[int], first *int, before *entgql.Cursor[int], last *int, where *ent.SubjectWhereInput) (*ent.SubjectConnection, error) {
	return r.client.Subject.Query().Paginate(ctx, after, first, before, last, ent.WithSubjectFilter(where.Filter))
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
