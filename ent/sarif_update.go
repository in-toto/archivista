// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/in-toto/archivista/ent/predicate"
	"github.com/in-toto/archivista/ent/sarif"
	"github.com/in-toto/archivista/ent/sarifrule"
	"github.com/in-toto/archivista/ent/statement"
)

// SarifUpdate is the builder for updating Sarif entities.
type SarifUpdate struct {
	config
	hooks    []Hook
	mutation *SarifMutation
}

// Where appends a list predicates to the SarifUpdate builder.
func (su *SarifUpdate) Where(ps ...predicate.Sarif) *SarifUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetReportFileName sets the "report_file_name" field.
func (su *SarifUpdate) SetReportFileName(s string) *SarifUpdate {
	su.mutation.SetReportFileName(s)
	return su
}

// SetNillableReportFileName sets the "report_file_name" field if the given value is not nil.
func (su *SarifUpdate) SetNillableReportFileName(s *string) *SarifUpdate {
	if s != nil {
		su.SetReportFileName(*s)
	}
	return su
}

// SetSarifRulesID sets the "sarif_rules" edge to the SarifRule entity by ID.
func (su *SarifUpdate) SetSarifRulesID(id uuid.UUID) *SarifUpdate {
	su.mutation.SetSarifRulesID(id)
	return su
}

// SetNillableSarifRulesID sets the "sarif_rules" edge to the SarifRule entity by ID if the given value is not nil.
func (su *SarifUpdate) SetNillableSarifRulesID(id *uuid.UUID) *SarifUpdate {
	if id != nil {
		su = su.SetSarifRulesID(*id)
	}
	return su
}

// SetSarifRules sets the "sarif_rules" edge to the SarifRule entity.
func (su *SarifUpdate) SetSarifRules(s *SarifRule) *SarifUpdate {
	return su.SetSarifRulesID(s.ID)
}

// SetStatementID sets the "statement" edge to the Statement entity by ID.
func (su *SarifUpdate) SetStatementID(id uuid.UUID) *SarifUpdate {
	su.mutation.SetStatementID(id)
	return su
}

// SetStatement sets the "statement" edge to the Statement entity.
func (su *SarifUpdate) SetStatement(s *Statement) *SarifUpdate {
	return su.SetStatementID(s.ID)
}

// Mutation returns the SarifMutation object of the builder.
func (su *SarifUpdate) Mutation() *SarifMutation {
	return su.mutation
}

// ClearSarifRules clears the "sarif_rules" edge to the SarifRule entity.
func (su *SarifUpdate) ClearSarifRules() *SarifUpdate {
	su.mutation.ClearSarifRules()
	return su
}

// ClearStatement clears the "statement" edge to the Statement entity.
func (su *SarifUpdate) ClearStatement() *SarifUpdate {
	su.mutation.ClearStatement()
	return su
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SarifUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SarifUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SarifUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SarifUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SarifUpdate) check() error {
	if v, ok := su.mutation.ReportFileName(); ok {
		if err := sarif.ReportFileNameValidator(v); err != nil {
			return &ValidationError{Name: "report_file_name", err: fmt.Errorf(`ent: validator failed for field "Sarif.report_file_name": %w`, err)}
		}
	}
	if _, ok := su.mutation.StatementID(); su.mutation.StatementCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Sarif.statement"`)
	}
	return nil
}

func (su *SarifUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(sarif.Table, sarif.Columns, sqlgraph.NewFieldSpec(sarif.FieldID, field.TypeUUID))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.ReportFileName(); ok {
		_spec.SetField(sarif.FieldReportFileName, field.TypeString, value)
	}
	if su.mutation.SarifRulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   sarif.SarifRulesTable,
			Columns: []string{sarif.SarifRulesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sarifrule.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.SarifRulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   sarif.SarifRulesTable,
			Columns: []string{sarif.SarifRulesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sarifrule.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.StatementCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   sarif.StatementTable,
			Columns: []string{sarif.StatementColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statement.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.StatementIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   sarif.StatementTable,
			Columns: []string{sarif.StatementColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statement.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{sarif.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SarifUpdateOne is the builder for updating a single Sarif entity.
type SarifUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SarifMutation
}

// SetReportFileName sets the "report_file_name" field.
func (suo *SarifUpdateOne) SetReportFileName(s string) *SarifUpdateOne {
	suo.mutation.SetReportFileName(s)
	return suo
}

// SetNillableReportFileName sets the "report_file_name" field if the given value is not nil.
func (suo *SarifUpdateOne) SetNillableReportFileName(s *string) *SarifUpdateOne {
	if s != nil {
		suo.SetReportFileName(*s)
	}
	return suo
}

// SetSarifRulesID sets the "sarif_rules" edge to the SarifRule entity by ID.
func (suo *SarifUpdateOne) SetSarifRulesID(id uuid.UUID) *SarifUpdateOne {
	suo.mutation.SetSarifRulesID(id)
	return suo
}

// SetNillableSarifRulesID sets the "sarif_rules" edge to the SarifRule entity by ID if the given value is not nil.
func (suo *SarifUpdateOne) SetNillableSarifRulesID(id *uuid.UUID) *SarifUpdateOne {
	if id != nil {
		suo = suo.SetSarifRulesID(*id)
	}
	return suo
}

// SetSarifRules sets the "sarif_rules" edge to the SarifRule entity.
func (suo *SarifUpdateOne) SetSarifRules(s *SarifRule) *SarifUpdateOne {
	return suo.SetSarifRulesID(s.ID)
}

// SetStatementID sets the "statement" edge to the Statement entity by ID.
func (suo *SarifUpdateOne) SetStatementID(id uuid.UUID) *SarifUpdateOne {
	suo.mutation.SetStatementID(id)
	return suo
}

// SetStatement sets the "statement" edge to the Statement entity.
func (suo *SarifUpdateOne) SetStatement(s *Statement) *SarifUpdateOne {
	return suo.SetStatementID(s.ID)
}

// Mutation returns the SarifMutation object of the builder.
func (suo *SarifUpdateOne) Mutation() *SarifMutation {
	return suo.mutation
}

// ClearSarifRules clears the "sarif_rules" edge to the SarifRule entity.
func (suo *SarifUpdateOne) ClearSarifRules() *SarifUpdateOne {
	suo.mutation.ClearSarifRules()
	return suo
}

// ClearStatement clears the "statement" edge to the Statement entity.
func (suo *SarifUpdateOne) ClearStatement() *SarifUpdateOne {
	suo.mutation.ClearStatement()
	return suo
}

// Where appends a list predicates to the SarifUpdate builder.
func (suo *SarifUpdateOne) Where(ps ...predicate.Sarif) *SarifUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SarifUpdateOne) Select(field string, fields ...string) *SarifUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Sarif entity.
func (suo *SarifUpdateOne) Save(ctx context.Context) (*Sarif, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SarifUpdateOne) SaveX(ctx context.Context) *Sarif {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SarifUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SarifUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SarifUpdateOne) check() error {
	if v, ok := suo.mutation.ReportFileName(); ok {
		if err := sarif.ReportFileNameValidator(v); err != nil {
			return &ValidationError{Name: "report_file_name", err: fmt.Errorf(`ent: validator failed for field "Sarif.report_file_name": %w`, err)}
		}
	}
	if _, ok := suo.mutation.StatementID(); suo.mutation.StatementCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "Sarif.statement"`)
	}
	return nil
}

func (suo *SarifUpdateOne) sqlSave(ctx context.Context) (_node *Sarif, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(sarif.Table, sarif.Columns, sqlgraph.NewFieldSpec(sarif.FieldID, field.TypeUUID))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Sarif.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, sarif.FieldID)
		for _, f := range fields {
			if !sarif.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != sarif.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := suo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := suo.mutation.ReportFileName(); ok {
		_spec.SetField(sarif.FieldReportFileName, field.TypeString, value)
	}
	if suo.mutation.SarifRulesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   sarif.SarifRulesTable,
			Columns: []string{sarif.SarifRulesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sarifrule.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.SarifRulesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   sarif.SarifRulesTable,
			Columns: []string{sarif.SarifRulesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(sarifrule.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.StatementCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   sarif.StatementTable,
			Columns: []string{sarif.StatementColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statement.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.StatementIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   sarif.StatementTable,
			Columns: []string{sarif.StatementColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(statement.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Sarif{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{sarif.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
