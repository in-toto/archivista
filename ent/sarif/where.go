// Code generated by ent, DO NOT EDIT.

package sarif

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/in-toto/archivista/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.Sarif {
	return predicate.Sarif(sql.FieldLTE(FieldID, id))
}

// ReportFileName applies equality check predicate on the "report_file_name" field. It's identical to ReportFileNameEQ.
func ReportFileName(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldEQ(FieldReportFileName, v))
}

// ReportFileNameEQ applies the EQ predicate on the "report_file_name" field.
func ReportFileNameEQ(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldEQ(FieldReportFileName, v))
}

// ReportFileNameNEQ applies the NEQ predicate on the "report_file_name" field.
func ReportFileNameNEQ(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldNEQ(FieldReportFileName, v))
}

// ReportFileNameIn applies the In predicate on the "report_file_name" field.
func ReportFileNameIn(vs ...string) predicate.Sarif {
	return predicate.Sarif(sql.FieldIn(FieldReportFileName, vs...))
}

// ReportFileNameNotIn applies the NotIn predicate on the "report_file_name" field.
func ReportFileNameNotIn(vs ...string) predicate.Sarif {
	return predicate.Sarif(sql.FieldNotIn(FieldReportFileName, vs...))
}

// ReportFileNameGT applies the GT predicate on the "report_file_name" field.
func ReportFileNameGT(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldGT(FieldReportFileName, v))
}

// ReportFileNameGTE applies the GTE predicate on the "report_file_name" field.
func ReportFileNameGTE(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldGTE(FieldReportFileName, v))
}

// ReportFileNameLT applies the LT predicate on the "report_file_name" field.
func ReportFileNameLT(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldLT(FieldReportFileName, v))
}

// ReportFileNameLTE applies the LTE predicate on the "report_file_name" field.
func ReportFileNameLTE(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldLTE(FieldReportFileName, v))
}

// ReportFileNameContains applies the Contains predicate on the "report_file_name" field.
func ReportFileNameContains(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldContains(FieldReportFileName, v))
}

// ReportFileNameHasPrefix applies the HasPrefix predicate on the "report_file_name" field.
func ReportFileNameHasPrefix(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldHasPrefix(FieldReportFileName, v))
}

// ReportFileNameHasSuffix applies the HasSuffix predicate on the "report_file_name" field.
func ReportFileNameHasSuffix(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldHasSuffix(FieldReportFileName, v))
}

// ReportFileNameEqualFold applies the EqualFold predicate on the "report_file_name" field.
func ReportFileNameEqualFold(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldEqualFold(FieldReportFileName, v))
}

// ReportFileNameContainsFold applies the ContainsFold predicate on the "report_file_name" field.
func ReportFileNameContainsFold(v string) predicate.Sarif {
	return predicate.Sarif(sql.FieldContainsFold(FieldReportFileName, v))
}

// HasSarifRules applies the HasEdge predicate on the "sarif_rules" edge.
func HasSarifRules() predicate.Sarif {
	return predicate.Sarif(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, SarifRulesTable, SarifRulesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSarifRulesWith applies the HasEdge predicate on the "sarif_rules" edge with a given conditions (other predicates).
func HasSarifRulesWith(preds ...predicate.SarifRule) predicate.Sarif {
	return predicate.Sarif(func(s *sql.Selector) {
		step := newSarifRulesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasStatement applies the HasEdge predicate on the "statement" edge.
func HasStatement() predicate.Sarif {
	return predicate.Sarif(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, StatementTable, StatementColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasStatementWith applies the HasEdge predicate on the "statement" edge with a given conditions (other predicates).
func HasStatementWith(preds ...predicate.Statement) predicate.Sarif {
	return predicate.Sarif(func(s *sql.Selector) {
		step := newStatementStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Sarif) predicate.Sarif {
	return predicate.Sarif(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Sarif) predicate.Sarif {
	return predicate.Sarif(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Sarif) predicate.Sarif {
	return predicate.Sarif(sql.NotPredicates(p))
}
