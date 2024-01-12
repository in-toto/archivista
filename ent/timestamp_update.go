// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/in-toto/archivista/ent/predicate"
	"github.com/in-toto/archivista/ent/signature"
	"github.com/in-toto/archivista/ent/timestamp"
)

// TimestampUpdate is the builder for updating Timestamp entities.
type TimestampUpdate struct {
	config
	hooks    []Hook
	mutation *TimestampMutation
}

// Where appends a list predicates to the TimestampUpdate builder.
func (tu *TimestampUpdate) Where(ps ...predicate.Timestamp) *TimestampUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetType sets the "type" field.
func (tu *TimestampUpdate) SetType(s string) *TimestampUpdate {
	tu.mutation.SetType(s)
	return tu
}

// SetNillableType sets the "type" field if the given value is not nil.
func (tu *TimestampUpdate) SetNillableType(s *string) *TimestampUpdate {
	if s != nil {
		tu.SetType(*s)
	}
	return tu
}

// SetTimestamp sets the "timestamp" field.
func (tu *TimestampUpdate) SetTimestamp(t time.Time) *TimestampUpdate {
	tu.mutation.SetTimestamp(t)
	return tu
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (tu *TimestampUpdate) SetNillableTimestamp(t *time.Time) *TimestampUpdate {
	if t != nil {
		tu.SetTimestamp(*t)
	}
	return tu
}

// SetSignatureID sets the "signature" edge to the Signature entity by ID.
func (tu *TimestampUpdate) SetSignatureID(id int) *TimestampUpdate {
	tu.mutation.SetSignatureID(id)
	return tu
}

// SetNillableSignatureID sets the "signature" edge to the Signature entity by ID if the given value is not nil.
func (tu *TimestampUpdate) SetNillableSignatureID(id *int) *TimestampUpdate {
	if id != nil {
		tu = tu.SetSignatureID(*id)
	}
	return tu
}

// SetSignature sets the "signature" edge to the Signature entity.
func (tu *TimestampUpdate) SetSignature(s *Signature) *TimestampUpdate {
	return tu.SetSignatureID(s.ID)
}

// Mutation returns the TimestampMutation object of the builder.
func (tu *TimestampUpdate) Mutation() *TimestampMutation {
	return tu.mutation
}

// ClearSignature clears the "signature" edge to the Signature entity.
func (tu *TimestampUpdate) ClearSignature() *TimestampUpdate {
	tu.mutation.ClearSignature()
	return tu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TimestampUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TimestampUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TimestampUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TimestampUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tu *TimestampUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(timestamp.Table, timestamp.Columns, sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.GetType(); ok {
		_spec.SetField(timestamp.FieldType, field.TypeString, value)
	}
	if value, ok := tu.mutation.Timestamp(); ok {
		_spec.SetField(timestamp.FieldTimestamp, field.TypeTime, value)
	}
	if tu.mutation.SignatureCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   timestamp.SignatureTable,
			Columns: []string{timestamp.SignatureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(signature.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.SignatureIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   timestamp.SignatureTable,
			Columns: []string{timestamp.SignatureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(signature.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{timestamp.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TimestampUpdateOne is the builder for updating a single Timestamp entity.
type TimestampUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TimestampMutation
}

// SetType sets the "type" field.
func (tuo *TimestampUpdateOne) SetType(s string) *TimestampUpdateOne {
	tuo.mutation.SetType(s)
	return tuo
}

// SetNillableType sets the "type" field if the given value is not nil.
func (tuo *TimestampUpdateOne) SetNillableType(s *string) *TimestampUpdateOne {
	if s != nil {
		tuo.SetType(*s)
	}
	return tuo
}

// SetTimestamp sets the "timestamp" field.
func (tuo *TimestampUpdateOne) SetTimestamp(t time.Time) *TimestampUpdateOne {
	tuo.mutation.SetTimestamp(t)
	return tuo
}

// SetNillableTimestamp sets the "timestamp" field if the given value is not nil.
func (tuo *TimestampUpdateOne) SetNillableTimestamp(t *time.Time) *TimestampUpdateOne {
	if t != nil {
		tuo.SetTimestamp(*t)
	}
	return tuo
}

// SetSignatureID sets the "signature" edge to the Signature entity by ID.
func (tuo *TimestampUpdateOne) SetSignatureID(id int) *TimestampUpdateOne {
	tuo.mutation.SetSignatureID(id)
	return tuo
}

// SetNillableSignatureID sets the "signature" edge to the Signature entity by ID if the given value is not nil.
func (tuo *TimestampUpdateOne) SetNillableSignatureID(id *int) *TimestampUpdateOne {
	if id != nil {
		tuo = tuo.SetSignatureID(*id)
	}
	return tuo
}

// SetSignature sets the "signature" edge to the Signature entity.
func (tuo *TimestampUpdateOne) SetSignature(s *Signature) *TimestampUpdateOne {
	return tuo.SetSignatureID(s.ID)
}

// Mutation returns the TimestampMutation object of the builder.
func (tuo *TimestampUpdateOne) Mutation() *TimestampMutation {
	return tuo.mutation
}

// ClearSignature clears the "signature" edge to the Signature entity.
func (tuo *TimestampUpdateOne) ClearSignature() *TimestampUpdateOne {
	tuo.mutation.ClearSignature()
	return tuo
}

// Where appends a list predicates to the TimestampUpdate builder.
func (tuo *TimestampUpdateOne) Where(ps ...predicate.Timestamp) *TimestampUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TimestampUpdateOne) Select(field string, fields ...string) *TimestampUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Timestamp entity.
func (tuo *TimestampUpdateOne) Save(ctx context.Context) (*Timestamp, error) {
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TimestampUpdateOne) SaveX(ctx context.Context) *Timestamp {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TimestampUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TimestampUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tuo *TimestampUpdateOne) sqlSave(ctx context.Context) (_node *Timestamp, err error) {
	_spec := sqlgraph.NewUpdateSpec(timestamp.Table, timestamp.Columns, sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Timestamp.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, timestamp.FieldID)
		for _, f := range fields {
			if !timestamp.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != timestamp.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.GetType(); ok {
		_spec.SetField(timestamp.FieldType, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Timestamp(); ok {
		_spec.SetField(timestamp.FieldTimestamp, field.TypeTime, value)
	}
	if tuo.mutation.SignatureCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   timestamp.SignatureTable,
			Columns: []string{timestamp.SignatureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(signature.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.SignatureIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   timestamp.SignatureTable,
			Columns: []string{timestamp.SignatureColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(signature.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Timestamp{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{timestamp.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}
