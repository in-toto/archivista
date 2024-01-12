// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/in-toto/archivista/ent/attestationcollection"
	"github.com/in-toto/archivista/ent/dsse"
	"github.com/in-toto/archivista/ent/predicate"
	"github.com/in-toto/archivista/ent/statement"
	"github.com/in-toto/archivista/ent/subject"
)

// StatementUpdate is the builder for updating Statement entities.
type StatementUpdate struct {
	config
	hooks    []Hook
	mutation *StatementMutation
}

// Where appends a list predicates to the StatementUpdate builder.
func (su *StatementUpdate) Where(ps ...predicate.Statement) *StatementUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetPredicate sets the "predicate" field.
func (su *StatementUpdate) SetPredicate(s string) *StatementUpdate {
	su.mutation.SetPredicate(s)
	return su
}

// SetNillablePredicate sets the "predicate" field if the given value is not nil.
func (su *StatementUpdate) SetNillablePredicate(s *string) *StatementUpdate {
	if s != nil {
		su.SetPredicate(*s)
	}
	return su
}

// AddSubjectIDs adds the "subjects" edge to the Subject entity by IDs.
func (su *StatementUpdate) AddSubjectIDs(ids ...int) *StatementUpdate {
	su.mutation.AddSubjectIDs(ids...)
	return su
}

// AddSubjects adds the "subjects" edges to the Subject entity.
func (su *StatementUpdate) AddSubjects(s ...*Subject) *StatementUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.AddSubjectIDs(ids...)
}

// SetAttestationCollectionsID sets the "attestation_collections" edge to the AttestationCollection entity by ID.
func (su *StatementUpdate) SetAttestationCollectionsID(id int) *StatementUpdate {
	su.mutation.SetAttestationCollectionsID(id)
	return su
}

// SetNillableAttestationCollectionsID sets the "attestation_collections" edge to the AttestationCollection entity by ID if the given value is not nil.
func (su *StatementUpdate) SetNillableAttestationCollectionsID(id *int) *StatementUpdate {
	if id != nil {
		su = su.SetAttestationCollectionsID(*id)
	}
	return su
}

// SetAttestationCollections sets the "attestation_collections" edge to the AttestationCollection entity.
func (su *StatementUpdate) SetAttestationCollections(a *AttestationCollection) *StatementUpdate {
	return su.SetAttestationCollectionsID(a.ID)
}

// AddDsseIDs adds the "dsse" edge to the Dsse entity by IDs.
func (su *StatementUpdate) AddDsseIDs(ids ...int) *StatementUpdate {
	su.mutation.AddDsseIDs(ids...)
	return su
}

// AddDsse adds the "dsse" edges to the Dsse entity.
func (su *StatementUpdate) AddDsse(d ...*Dsse) *StatementUpdate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return su.AddDsseIDs(ids...)
}

// Mutation returns the StatementMutation object of the builder.
func (su *StatementUpdate) Mutation() *StatementMutation {
	return su.mutation
}

// ClearSubjects clears all "subjects" edges to the Subject entity.
func (su *StatementUpdate) ClearSubjects() *StatementUpdate {
	su.mutation.ClearSubjects()
	return su
}

// RemoveSubjectIDs removes the "subjects" edge to Subject entities by IDs.
func (su *StatementUpdate) RemoveSubjectIDs(ids ...int) *StatementUpdate {
	su.mutation.RemoveSubjectIDs(ids...)
	return su
}

// RemoveSubjects removes "subjects" edges to Subject entities.
func (su *StatementUpdate) RemoveSubjects(s ...*Subject) *StatementUpdate {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return su.RemoveSubjectIDs(ids...)
}

// ClearAttestationCollections clears the "attestation_collections" edge to the AttestationCollection entity.
func (su *StatementUpdate) ClearAttestationCollections() *StatementUpdate {
	su.mutation.ClearAttestationCollections()
	return su
}

// ClearDsse clears all "dsse" edges to the Dsse entity.
func (su *StatementUpdate) ClearDsse() *StatementUpdate {
	su.mutation.ClearDsse()
	return su
}

// RemoveDsseIDs removes the "dsse" edge to Dsse entities by IDs.
func (su *StatementUpdate) RemoveDsseIDs(ids ...int) *StatementUpdate {
	su.mutation.RemoveDsseIDs(ids...)
	return su
}

// RemoveDsse removes "dsse" edges to Dsse entities.
func (su *StatementUpdate) RemoveDsse(d ...*Dsse) *StatementUpdate {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return su.RemoveDsseIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *StatementUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *StatementUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *StatementUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *StatementUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *StatementUpdate) check() error {
	if v, ok := su.mutation.Predicate(); ok {
		if err := statement.PredicateValidator(v); err != nil {
			return &ValidationError{Name: "predicate", err: fmt.Errorf(`ent: validator failed for field "Statement.predicate": %w`, err)}
		}
	}
	return nil
}

func (su *StatementUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(statement.Table, statement.Columns, sqlgraph.NewFieldSpec(statement.FieldID, field.TypeInt))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.Predicate(); ok {
		_spec.SetField(statement.FieldPredicate, field.TypeString, value)
	}
	if su.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubjectsTable,
			Columns: []string{statement.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subject.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedSubjectsIDs(); len(nodes) > 0 && !su.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubjectsTable,
			Columns: []string{statement.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subject.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.SubjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubjectsTable,
			Columns: []string{statement.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subject.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.AttestationCollectionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   statement.AttestationCollectionsTable,
			Columns: []string{statement.AttestationCollectionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(attestationcollection.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.AttestationCollectionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   statement.AttestationCollectionsTable,
			Columns: []string{statement.AttestationCollectionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(attestationcollection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if su.mutation.DsseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   statement.DsseTable,
			Columns: []string{statement.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedDsseIDs(); len(nodes) > 0 && !su.mutation.DsseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   statement.DsseTable,
			Columns: []string{statement.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.DsseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   statement.DsseTable,
			Columns: []string{statement.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{statement.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// StatementUpdateOne is the builder for updating a single Statement entity.
type StatementUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *StatementMutation
}

// SetPredicate sets the "predicate" field.
func (suo *StatementUpdateOne) SetPredicate(s string) *StatementUpdateOne {
	suo.mutation.SetPredicate(s)
	return suo
}

// SetNillablePredicate sets the "predicate" field if the given value is not nil.
func (suo *StatementUpdateOne) SetNillablePredicate(s *string) *StatementUpdateOne {
	if s != nil {
		suo.SetPredicate(*s)
	}
	return suo
}

// AddSubjectIDs adds the "subjects" edge to the Subject entity by IDs.
func (suo *StatementUpdateOne) AddSubjectIDs(ids ...int) *StatementUpdateOne {
	suo.mutation.AddSubjectIDs(ids...)
	return suo
}

// AddSubjects adds the "subjects" edges to the Subject entity.
func (suo *StatementUpdateOne) AddSubjects(s ...*Subject) *StatementUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.AddSubjectIDs(ids...)
}

// SetAttestationCollectionsID sets the "attestation_collections" edge to the AttestationCollection entity by ID.
func (suo *StatementUpdateOne) SetAttestationCollectionsID(id int) *StatementUpdateOne {
	suo.mutation.SetAttestationCollectionsID(id)
	return suo
}

// SetNillableAttestationCollectionsID sets the "attestation_collections" edge to the AttestationCollection entity by ID if the given value is not nil.
func (suo *StatementUpdateOne) SetNillableAttestationCollectionsID(id *int) *StatementUpdateOne {
	if id != nil {
		suo = suo.SetAttestationCollectionsID(*id)
	}
	return suo
}

// SetAttestationCollections sets the "attestation_collections" edge to the AttestationCollection entity.
func (suo *StatementUpdateOne) SetAttestationCollections(a *AttestationCollection) *StatementUpdateOne {
	return suo.SetAttestationCollectionsID(a.ID)
}

// AddDsseIDs adds the "dsse" edge to the Dsse entity by IDs.
func (suo *StatementUpdateOne) AddDsseIDs(ids ...int) *StatementUpdateOne {
	suo.mutation.AddDsseIDs(ids...)
	return suo
}

// AddDsse adds the "dsse" edges to the Dsse entity.
func (suo *StatementUpdateOne) AddDsse(d ...*Dsse) *StatementUpdateOne {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return suo.AddDsseIDs(ids...)
}

// Mutation returns the StatementMutation object of the builder.
func (suo *StatementUpdateOne) Mutation() *StatementMutation {
	return suo.mutation
}

// ClearSubjects clears all "subjects" edges to the Subject entity.
func (suo *StatementUpdateOne) ClearSubjects() *StatementUpdateOne {
	suo.mutation.ClearSubjects()
	return suo
}

// RemoveSubjectIDs removes the "subjects" edge to Subject entities by IDs.
func (suo *StatementUpdateOne) RemoveSubjectIDs(ids ...int) *StatementUpdateOne {
	suo.mutation.RemoveSubjectIDs(ids...)
	return suo
}

// RemoveSubjects removes "subjects" edges to Subject entities.
func (suo *StatementUpdateOne) RemoveSubjects(s ...*Subject) *StatementUpdateOne {
	ids := make([]int, len(s))
	for i := range s {
		ids[i] = s[i].ID
	}
	return suo.RemoveSubjectIDs(ids...)
}

// ClearAttestationCollections clears the "attestation_collections" edge to the AttestationCollection entity.
func (suo *StatementUpdateOne) ClearAttestationCollections() *StatementUpdateOne {
	suo.mutation.ClearAttestationCollections()
	return suo
}

// ClearDsse clears all "dsse" edges to the Dsse entity.
func (suo *StatementUpdateOne) ClearDsse() *StatementUpdateOne {
	suo.mutation.ClearDsse()
	return suo
}

// RemoveDsseIDs removes the "dsse" edge to Dsse entities by IDs.
func (suo *StatementUpdateOne) RemoveDsseIDs(ids ...int) *StatementUpdateOne {
	suo.mutation.RemoveDsseIDs(ids...)
	return suo
}

// RemoveDsse removes "dsse" edges to Dsse entities.
func (suo *StatementUpdateOne) RemoveDsse(d ...*Dsse) *StatementUpdateOne {
	ids := make([]int, len(d))
	for i := range d {
		ids[i] = d[i].ID
	}
	return suo.RemoveDsseIDs(ids...)
}

// Where appends a list predicates to the StatementUpdate builder.
func (suo *StatementUpdateOne) Where(ps ...predicate.Statement) *StatementUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *StatementUpdateOne) Select(field string, fields ...string) *StatementUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Statement entity.
func (suo *StatementUpdateOne) Save(ctx context.Context) (*Statement, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *StatementUpdateOne) SaveX(ctx context.Context) *Statement {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *StatementUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *StatementUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *StatementUpdateOne) check() error {
	if v, ok := suo.mutation.Predicate(); ok {
		if err := statement.PredicateValidator(v); err != nil {
			return &ValidationError{Name: "predicate", err: fmt.Errorf(`ent: validator failed for field "Statement.predicate": %w`, err)}
		}
	}
	return nil
}

func (suo *StatementUpdateOne) sqlSave(ctx context.Context) (_node *Statement, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(statement.Table, statement.Columns, sqlgraph.NewFieldSpec(statement.FieldID, field.TypeInt))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Statement.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, statement.FieldID)
		for _, f := range fields {
			if !statement.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != statement.FieldID {
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
	if value, ok := suo.mutation.Predicate(); ok {
		_spec.SetField(statement.FieldPredicate, field.TypeString, value)
	}
	if suo.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubjectsTable,
			Columns: []string{statement.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subject.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedSubjectsIDs(); len(nodes) > 0 && !suo.mutation.SubjectsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubjectsTable,
			Columns: []string{statement.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subject.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.SubjectsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   statement.SubjectsTable,
			Columns: []string{statement.SubjectsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(subject.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.AttestationCollectionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   statement.AttestationCollectionsTable,
			Columns: []string{statement.AttestationCollectionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(attestationcollection.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.AttestationCollectionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   statement.AttestationCollectionsTable,
			Columns: []string{statement.AttestationCollectionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(attestationcollection.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if suo.mutation.DsseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   statement.DsseTable,
			Columns: []string{statement.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedDsseIDs(); len(nodes) > 0 && !suo.mutation.DsseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   statement.DsseTable,
			Columns: []string{statement.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.DsseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: true,
			Table:   statement.DsseTable,
			Columns: []string{statement.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Statement{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{statement.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
