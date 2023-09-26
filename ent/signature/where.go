// Code generated by ent, DO NOT EDIT.

package signature

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/testifysec/archivista/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Signature {
	return predicate.Signature(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Signature {
	return predicate.Signature(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Signature {
	return predicate.Signature(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Signature {
	return predicate.Signature(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Signature {
	return predicate.Signature(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Signature {
	return predicate.Signature(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Signature {
	return predicate.Signature(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Signature {
	return predicate.Signature(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Signature {
	return predicate.Signature(sql.FieldLTE(FieldID, id))
}

// KeyID applies equality check predicate on the "key_id" field. It's identical to KeyIDEQ.
func KeyID(v string) predicate.Signature {
	return predicate.Signature(sql.FieldEQ(FieldKeyID, v))
}

// Signature applies equality check predicate on the "signature" field. It's identical to SignatureEQ.
func Signature(v string) predicate.Signature {
	return predicate.Signature(sql.FieldEQ(FieldSignature, v))
}

// KeyIDEQ applies the EQ predicate on the "key_id" field.
func KeyIDEQ(v string) predicate.Signature {
	return predicate.Signature(sql.FieldEQ(FieldKeyID, v))
}

// KeyIDNEQ applies the NEQ predicate on the "key_id" field.
func KeyIDNEQ(v string) predicate.Signature {
	return predicate.Signature(sql.FieldNEQ(FieldKeyID, v))
}

// KeyIDIn applies the In predicate on the "key_id" field.
func KeyIDIn(vs ...string) predicate.Signature {
	return predicate.Signature(sql.FieldIn(FieldKeyID, vs...))
}

// KeyIDNotIn applies the NotIn predicate on the "key_id" field.
func KeyIDNotIn(vs ...string) predicate.Signature {
	return predicate.Signature(sql.FieldNotIn(FieldKeyID, vs...))
}

// KeyIDGT applies the GT predicate on the "key_id" field.
func KeyIDGT(v string) predicate.Signature {
	return predicate.Signature(sql.FieldGT(FieldKeyID, v))
}

// KeyIDGTE applies the GTE predicate on the "key_id" field.
func KeyIDGTE(v string) predicate.Signature {
	return predicate.Signature(sql.FieldGTE(FieldKeyID, v))
}

// KeyIDLT applies the LT predicate on the "key_id" field.
func KeyIDLT(v string) predicate.Signature {
	return predicate.Signature(sql.FieldLT(FieldKeyID, v))
}

// KeyIDLTE applies the LTE predicate on the "key_id" field.
func KeyIDLTE(v string) predicate.Signature {
	return predicate.Signature(sql.FieldLTE(FieldKeyID, v))
}

// KeyIDContains applies the Contains predicate on the "key_id" field.
func KeyIDContains(v string) predicate.Signature {
	return predicate.Signature(sql.FieldContains(FieldKeyID, v))
}

// KeyIDHasPrefix applies the HasPrefix predicate on the "key_id" field.
func KeyIDHasPrefix(v string) predicate.Signature {
	return predicate.Signature(sql.FieldHasPrefix(FieldKeyID, v))
}

// KeyIDHasSuffix applies the HasSuffix predicate on the "key_id" field.
func KeyIDHasSuffix(v string) predicate.Signature {
	return predicate.Signature(sql.FieldHasSuffix(FieldKeyID, v))
}

// KeyIDEqualFold applies the EqualFold predicate on the "key_id" field.
func KeyIDEqualFold(v string) predicate.Signature {
	return predicate.Signature(sql.FieldEqualFold(FieldKeyID, v))
}

// KeyIDContainsFold applies the ContainsFold predicate on the "key_id" field.
func KeyIDContainsFold(v string) predicate.Signature {
	return predicate.Signature(sql.FieldContainsFold(FieldKeyID, v))
}

// SignatureEQ applies the EQ predicate on the "signature" field.
func SignatureEQ(v string) predicate.Signature {
	return predicate.Signature(sql.FieldEQ(FieldSignature, v))
}

// SignatureNEQ applies the NEQ predicate on the "signature" field.
func SignatureNEQ(v string) predicate.Signature {
	return predicate.Signature(sql.FieldNEQ(FieldSignature, v))
}

// SignatureIn applies the In predicate on the "signature" field.
func SignatureIn(vs ...string) predicate.Signature {
	return predicate.Signature(sql.FieldIn(FieldSignature, vs...))
}

// SignatureNotIn applies the NotIn predicate on the "signature" field.
func SignatureNotIn(vs ...string) predicate.Signature {
	return predicate.Signature(sql.FieldNotIn(FieldSignature, vs...))
}

// SignatureGT applies the GT predicate on the "signature" field.
func SignatureGT(v string) predicate.Signature {
	return predicate.Signature(sql.FieldGT(FieldSignature, v))
}

// SignatureGTE applies the GTE predicate on the "signature" field.
func SignatureGTE(v string) predicate.Signature {
	return predicate.Signature(sql.FieldGTE(FieldSignature, v))
}

// SignatureLT applies the LT predicate on the "signature" field.
func SignatureLT(v string) predicate.Signature {
	return predicate.Signature(sql.FieldLT(FieldSignature, v))
}

// SignatureLTE applies the LTE predicate on the "signature" field.
func SignatureLTE(v string) predicate.Signature {
	return predicate.Signature(sql.FieldLTE(FieldSignature, v))
}

// SignatureContains applies the Contains predicate on the "signature" field.
func SignatureContains(v string) predicate.Signature {
	return predicate.Signature(sql.FieldContains(FieldSignature, v))
}

// SignatureHasPrefix applies the HasPrefix predicate on the "signature" field.
func SignatureHasPrefix(v string) predicate.Signature {
	return predicate.Signature(sql.FieldHasPrefix(FieldSignature, v))
}

// SignatureHasSuffix applies the HasSuffix predicate on the "signature" field.
func SignatureHasSuffix(v string) predicate.Signature {
	return predicate.Signature(sql.FieldHasSuffix(FieldSignature, v))
}

// SignatureEqualFold applies the EqualFold predicate on the "signature" field.
func SignatureEqualFold(v string) predicate.Signature {
	return predicate.Signature(sql.FieldEqualFold(FieldSignature, v))
}

// SignatureContainsFold applies the ContainsFold predicate on the "signature" field.
func SignatureContainsFold(v string) predicate.Signature {
	return predicate.Signature(sql.FieldContainsFold(FieldSignature, v))
}

// HasDsse applies the HasEdge predicate on the "dsse" edge.
func HasDsse() predicate.Signature {
	return predicate.Signature(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, DsseTable, DsseColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasDsseWith applies the HasEdge predicate on the "dsse" edge with a given conditions (other predicates).
func HasDsseWith(preds ...predicate.Dsse) predicate.Signature {
	return predicate.Signature(func(s *sql.Selector) {
		step := newDsseStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasTimestamps applies the HasEdge predicate on the "timestamps" edge.
func HasTimestamps() predicate.Signature {
	return predicate.Signature(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, TimestampsTable, TimestampsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTimestampsWith applies the HasEdge predicate on the "timestamps" edge with a given conditions (other predicates).
func HasTimestampsWith(preds ...predicate.Timestamp) predicate.Signature {
	return predicate.Signature(func(s *sql.Selector) {
		step := newTimestampsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Signature) predicate.Signature {
	return predicate.Signature(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Signature) predicate.Signature {
	return predicate.Signature(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Signature) predicate.Signature {
	return predicate.Signature(sql.NotPredicates(p))
}
