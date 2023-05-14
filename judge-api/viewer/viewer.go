package viewer

import (
	"context"

	"github.com/google/uuid"
)

// Viewer describes the query/mutation viewer-context.
type Viewer interface {
	Tenants() []uuid.UUID   // The user's assigned tenants
	UserIdentityID() string // The user's identity ID
}

type UserViewer struct {
	TenantsAssigned []uuid.UUID
	IdentityID      string
}

func (v UserViewer) Tenants() []uuid.UUID {
	if v.TenantsAssigned == nil {
		return []uuid.UUID{}
	}
	return v.TenantsAssigned
}

func (v UserViewer) UserIdentityID() string {
	return v.IdentityID
}

type ctxKey struct{}

// FromContext returns the identity stored in a context.
func FromContext(ctx context.Context) Viewer {
	v, _ := ctx.Value(ctxKey{}).(Viewer)
	return v
}

// NewContext returns a copy of parent context with the given Viewer attached to it.
func NewContext(parent context.Context, v Viewer) context.Context {
	return context.WithValue(parent, ctxKey{}, v)
}
