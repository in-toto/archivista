// Code generated by ent, DO NOT EDIT.

package dsse

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/testifysec/archivist/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		v := make([]any, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// GitoidSha256 applies equality check predicate on the "gitoid_sha256" field. It's identical to GitoidSha256EQ.
func GitoidSha256(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldGitoidSha256), v))
	})
}

// PayloadType applies equality check predicate on the "payload_type" field. It's identical to PayloadTypeEQ.
func PayloadType(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPayloadType), v))
	})
}

// GitoidSha256EQ applies the EQ predicate on the "gitoid_sha256" field.
func GitoidSha256EQ(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256NEQ applies the NEQ predicate on the "gitoid_sha256" field.
func GitoidSha256NEQ(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256In applies the In predicate on the "gitoid_sha256" field.
func GitoidSha256In(vs ...string) predicate.Dsse {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldGitoidSha256), v...))
	})
}

// GitoidSha256NotIn applies the NotIn predicate on the "gitoid_sha256" field.
func GitoidSha256NotIn(vs ...string) predicate.Dsse {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldGitoidSha256), v...))
	})
}

// GitoidSha256GT applies the GT predicate on the "gitoid_sha256" field.
func GitoidSha256GT(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256GTE applies the GTE predicate on the "gitoid_sha256" field.
func GitoidSha256GTE(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256LT applies the LT predicate on the "gitoid_sha256" field.
func GitoidSha256LT(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256LTE applies the LTE predicate on the "gitoid_sha256" field.
func GitoidSha256LTE(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256Contains applies the Contains predicate on the "gitoid_sha256" field.
func GitoidSha256Contains(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256HasPrefix applies the HasPrefix predicate on the "gitoid_sha256" field.
func GitoidSha256HasPrefix(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256HasSuffix applies the HasSuffix predicate on the "gitoid_sha256" field.
func GitoidSha256HasSuffix(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256EqualFold applies the EqualFold predicate on the "gitoid_sha256" field.
func GitoidSha256EqualFold(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldGitoidSha256), v))
	})
}

// GitoidSha256ContainsFold applies the ContainsFold predicate on the "gitoid_sha256" field.
func GitoidSha256ContainsFold(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldGitoidSha256), v))
	})
}

// PayloadTypeEQ applies the EQ predicate on the "payload_type" field.
func PayloadTypeEQ(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeNEQ applies the NEQ predicate on the "payload_type" field.
func PayloadTypeNEQ(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeIn applies the In predicate on the "payload_type" field.
func PayloadTypeIn(vs ...string) predicate.Dsse {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldPayloadType), v...))
	})
}

// PayloadTypeNotIn applies the NotIn predicate on the "payload_type" field.
func PayloadTypeNotIn(vs ...string) predicate.Dsse {
	v := make([]any, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldPayloadType), v...))
	})
}

// PayloadTypeGT applies the GT predicate on the "payload_type" field.
func PayloadTypeGT(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeGTE applies the GTE predicate on the "payload_type" field.
func PayloadTypeGTE(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeLT applies the LT predicate on the "payload_type" field.
func PayloadTypeLT(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeLTE applies the LTE predicate on the "payload_type" field.
func PayloadTypeLTE(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeContains applies the Contains predicate on the "payload_type" field.
func PayloadTypeContains(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeHasPrefix applies the HasPrefix predicate on the "payload_type" field.
func PayloadTypeHasPrefix(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeHasSuffix applies the HasSuffix predicate on the "payload_type" field.
func PayloadTypeHasSuffix(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeEqualFold applies the EqualFold predicate on the "payload_type" field.
func PayloadTypeEqualFold(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldPayloadType), v))
	})
}

// PayloadTypeContainsFold applies the ContainsFold predicate on the "payload_type" field.
func PayloadTypeContainsFold(v string) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldPayloadType), v))
	})
}

// HasStatement applies the HasEdge predicate on the "statement" edge.
func HasStatement() predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StatementTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, StatementTable, StatementColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStatementWith applies the HasEdge predicate on the "statement" edge with a given conditions (other predicates).
func HasStatementWith(preds ...predicate.Statement) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(StatementInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, StatementTable, StatementColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasSignatures applies the HasEdge predicate on the "signatures" edge.
func HasSignatures() predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SignaturesTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SignaturesTable, SignaturesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSignaturesWith applies the HasEdge predicate on the "signatures" edge with a given conditions (other predicates).
func HasSignaturesWith(preds ...predicate.Signature) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(SignaturesInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, SignaturesTable, SignaturesColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPayloadDigests applies the HasEdge predicate on the "payload_digests" edge.
func HasPayloadDigests() predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PayloadDigestsTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PayloadDigestsTable, PayloadDigestsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPayloadDigestsWith applies the HasEdge predicate on the "payload_digests" edge with a given conditions (other predicates).
func HasPayloadDigestsWith(preds ...predicate.PayloadDigest) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(PayloadDigestsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PayloadDigestsTable, PayloadDigestsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Dsse) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Dsse) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
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
func Not(p predicate.Dsse) predicate.Dsse {
	return predicate.Dsse(func(s *sql.Selector) {
		p(s.Not())
	})
}
