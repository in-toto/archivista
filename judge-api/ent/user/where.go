// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// IdentityID applies equality check predicate on the "identity_id" field. It's identical to IdentityIDEQ.
func IdentityID(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldIdentityID, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUpdatedAt, v))
}

// IdentityIDEQ applies the EQ predicate on the "identity_id" field.
func IdentityIDEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldIdentityID, v))
}

// IdentityIDNEQ applies the NEQ predicate on the "identity_id" field.
func IdentityIDNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldIdentityID, v))
}

// IdentityIDIn applies the In predicate on the "identity_id" field.
func IdentityIDIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldIdentityID, vs...))
}

// IdentityIDNotIn applies the NotIn predicate on the "identity_id" field.
func IdentityIDNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldIdentityID, vs...))
}

// IdentityIDGT applies the GT predicate on the "identity_id" field.
func IdentityIDGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldIdentityID, v))
}

// IdentityIDGTE applies the GTE predicate on the "identity_id" field.
func IdentityIDGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldIdentityID, v))
}

// IdentityIDLT applies the LT predicate on the "identity_id" field.
func IdentityIDLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldIdentityID, v))
}

// IdentityIDLTE applies the LTE predicate on the "identity_id" field.
func IdentityIDLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldIdentityID, v))
}

// IdentityIDContains applies the Contains predicate on the "identity_id" field.
func IdentityIDContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldIdentityID, v))
}

// IdentityIDHasPrefix applies the HasPrefix predicate on the "identity_id" field.
func IdentityIDHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldIdentityID, v))
}

// IdentityIDHasSuffix applies the HasSuffix predicate on the "identity_id" field.
func IdentityIDHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldIdentityID, v))
}

// IdentityIDEqualFold applies the EqualFold predicate on the "identity_id" field.
func IdentityIDEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldIdentityID, v))
}

// IdentityIDContainsFold applies the ContainsFold predicate on the "identity_id" field.
func IdentityIDContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldIdentityID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.User {
	return predicate.User(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.User {
	return predicate.User(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.User {
	return predicate.User(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.User {
	return predicate.User(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasTenant applies the HasEdge predicate on the "tenant" edge.
func HasTenant() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, TenantTable, TenantPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTenantWith applies the HasEdge predicate on the "tenant" edge with a given conditions (other predicates).
func HasTenantWith(preds ...predicate.Tenant) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(TenantInverseTable, TenantFieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, TenantTable, TenantPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.User) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		p(s.Not())
	})
}
