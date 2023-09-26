// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/testifysec/archivista/ent/dsse"
	"github.com/testifysec/archivista/ent/predicate"
	"github.com/testifysec/archivista/ent/signature"
	"github.com/testifysec/archivista/ent/timestamp"
)

// SignatureUpdate is the builder for updating Signature entities.
type SignatureUpdate struct {
	config
	hooks    []Hook
	mutation *SignatureMutation
}

// Where appends a list predicates to the SignatureUpdate builder.
func (su *SignatureUpdate) Where(ps ...predicate.Signature) *SignatureUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetKeyID sets the "key_id" field.
func (su *SignatureUpdate) SetKeyID(s string) *SignatureUpdate {
	su.mutation.SetKeyID(s)
	return su
}

// SetSignature sets the "signature" field.
func (su *SignatureUpdate) SetSignature(s string) *SignatureUpdate {
	su.mutation.SetSignature(s)
	return su
}

// SetDsseID sets the "dsse" edge to the Dsse entity by ID.
func (su *SignatureUpdate) SetDsseID(id int) *SignatureUpdate {
	su.mutation.SetDsseID(id)
	return su
}

// SetNillableDsseID sets the "dsse" edge to the Dsse entity by ID if the given value is not nil.
func (su *SignatureUpdate) SetNillableDsseID(id *int) *SignatureUpdate {
	if id != nil {
		su = su.SetDsseID(*id)
	}
	return su
}

// SetDsse sets the "dsse" edge to the Dsse entity.
func (su *SignatureUpdate) SetDsse(d *Dsse) *SignatureUpdate {
	return su.SetDsseID(d.ID)
}

// AddTimestampIDs adds the "timestamps" edge to the Timestamp entity by IDs.
func (su *SignatureUpdate) AddTimestampIDs(ids ...int) *SignatureUpdate {
	su.mutation.AddTimestampIDs(ids...)
	return su
}

// AddTimestamps adds the "timestamps" edges to the Timestamp entity.
func (su *SignatureUpdate) AddTimestamps(t ...*Timestamp) *SignatureUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return su.AddTimestampIDs(ids...)
}

// Mutation returns the SignatureMutation object of the builder.
func (su *SignatureUpdate) Mutation() *SignatureMutation {
	return su.mutation
}

// ClearDsse clears the "dsse" edge to the Dsse entity.
func (su *SignatureUpdate) ClearDsse() *SignatureUpdate {
	su.mutation.ClearDsse()
	return su
}

// ClearTimestamps clears all "timestamps" edges to the Timestamp entity.
func (su *SignatureUpdate) ClearTimestamps() *SignatureUpdate {
	su.mutation.ClearTimestamps()
	return su
}

// RemoveTimestampIDs removes the "timestamps" edge to Timestamp entities by IDs.
func (su *SignatureUpdate) RemoveTimestampIDs(ids ...int) *SignatureUpdate {
	su.mutation.RemoveTimestampIDs(ids...)
	return su
}

// RemoveTimestamps removes "timestamps" edges to Timestamp entities.
func (su *SignatureUpdate) RemoveTimestamps(t ...*Timestamp) *SignatureUpdate {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return su.RemoveTimestampIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SignatureUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, su.sqlSave, su.mutation, su.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (su *SignatureUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SignatureUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SignatureUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (su *SignatureUpdate) check() error {
	if v, ok := su.mutation.KeyID(); ok {
		if err := signature.KeyIDValidator(v); err != nil {
			return &ValidationError{Name: "key_id", err: fmt.Errorf(`ent: validator failed for field "Signature.key_id": %w`, err)}
		}
	}
	if v, ok := su.mutation.Signature(); ok {
		if err := signature.SignatureValidator(v); err != nil {
			return &ValidationError{Name: "signature", err: fmt.Errorf(`ent: validator failed for field "Signature.signature": %w`, err)}
		}
	}
	return nil
}

func (su *SignatureUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := su.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(signature.Table, signature.Columns, sqlgraph.NewFieldSpec(signature.FieldID, field.TypeInt))
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.KeyID(); ok {
		_spec.SetField(signature.FieldKeyID, field.TypeString, value)
	}
	if value, ok := su.mutation.Signature(); ok {
		_spec.SetField(signature.FieldSignature, field.TypeString, value)
	}
	if su.mutation.DsseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   signature.DsseTable,
			Columns: []string{signature.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.DsseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   signature.DsseTable,
			Columns: []string{signature.DsseColumn},
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
	if su.mutation.TimestampsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   signature.TimestampsTable,
			Columns: []string{signature.TimestampsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.RemovedTimestampsIDs(); len(nodes) > 0 && !su.mutation.TimestampsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   signature.TimestampsTable,
			Columns: []string{signature.TimestampsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := su.mutation.TimestampsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   signature.TimestampsTable,
			Columns: []string{signature.TimestampsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{signature.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	su.mutation.done = true
	return n, nil
}

// SignatureUpdateOne is the builder for updating a single Signature entity.
type SignatureUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SignatureMutation
}

// SetKeyID sets the "key_id" field.
func (suo *SignatureUpdateOne) SetKeyID(s string) *SignatureUpdateOne {
	suo.mutation.SetKeyID(s)
	return suo
}

// SetSignature sets the "signature" field.
func (suo *SignatureUpdateOne) SetSignature(s string) *SignatureUpdateOne {
	suo.mutation.SetSignature(s)
	return suo
}

// SetDsseID sets the "dsse" edge to the Dsse entity by ID.
func (suo *SignatureUpdateOne) SetDsseID(id int) *SignatureUpdateOne {
	suo.mutation.SetDsseID(id)
	return suo
}

// SetNillableDsseID sets the "dsse" edge to the Dsse entity by ID if the given value is not nil.
func (suo *SignatureUpdateOne) SetNillableDsseID(id *int) *SignatureUpdateOne {
	if id != nil {
		suo = suo.SetDsseID(*id)
	}
	return suo
}

// SetDsse sets the "dsse" edge to the Dsse entity.
func (suo *SignatureUpdateOne) SetDsse(d *Dsse) *SignatureUpdateOne {
	return suo.SetDsseID(d.ID)
}

// AddTimestampIDs adds the "timestamps" edge to the Timestamp entity by IDs.
func (suo *SignatureUpdateOne) AddTimestampIDs(ids ...int) *SignatureUpdateOne {
	suo.mutation.AddTimestampIDs(ids...)
	return suo
}

// AddTimestamps adds the "timestamps" edges to the Timestamp entity.
func (suo *SignatureUpdateOne) AddTimestamps(t ...*Timestamp) *SignatureUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return suo.AddTimestampIDs(ids...)
}

// Mutation returns the SignatureMutation object of the builder.
func (suo *SignatureUpdateOne) Mutation() *SignatureMutation {
	return suo.mutation
}

// ClearDsse clears the "dsse" edge to the Dsse entity.
func (suo *SignatureUpdateOne) ClearDsse() *SignatureUpdateOne {
	suo.mutation.ClearDsse()
	return suo
}

// ClearTimestamps clears all "timestamps" edges to the Timestamp entity.
func (suo *SignatureUpdateOne) ClearTimestamps() *SignatureUpdateOne {
	suo.mutation.ClearTimestamps()
	return suo
}

// RemoveTimestampIDs removes the "timestamps" edge to Timestamp entities by IDs.
func (suo *SignatureUpdateOne) RemoveTimestampIDs(ids ...int) *SignatureUpdateOne {
	suo.mutation.RemoveTimestampIDs(ids...)
	return suo
}

// RemoveTimestamps removes "timestamps" edges to Timestamp entities.
func (suo *SignatureUpdateOne) RemoveTimestamps(t ...*Timestamp) *SignatureUpdateOne {
	ids := make([]int, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return suo.RemoveTimestampIDs(ids...)
}

// Where appends a list predicates to the SignatureUpdate builder.
func (suo *SignatureUpdateOne) Where(ps ...predicate.Signature) *SignatureUpdateOne {
	suo.mutation.Where(ps...)
	return suo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SignatureUpdateOne) Select(field string, fields ...string) *SignatureUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Signature entity.
func (suo *SignatureUpdateOne) Save(ctx context.Context) (*Signature, error) {
	return withHooks(ctx, suo.sqlSave, suo.mutation, suo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SignatureUpdateOne) SaveX(ctx context.Context) *Signature {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SignatureUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SignatureUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (suo *SignatureUpdateOne) check() error {
	if v, ok := suo.mutation.KeyID(); ok {
		if err := signature.KeyIDValidator(v); err != nil {
			return &ValidationError{Name: "key_id", err: fmt.Errorf(`ent: validator failed for field "Signature.key_id": %w`, err)}
		}
	}
	if v, ok := suo.mutation.Signature(); ok {
		if err := signature.SignatureValidator(v); err != nil {
			return &ValidationError{Name: "signature", err: fmt.Errorf(`ent: validator failed for field "Signature.signature": %w`, err)}
		}
	}
	return nil
}

func (suo *SignatureUpdateOne) sqlSave(ctx context.Context) (_node *Signature, err error) {
	if err := suo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(signature.Table, signature.Columns, sqlgraph.NewFieldSpec(signature.FieldID, field.TypeInt))
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Signature.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, signature.FieldID)
		for _, f := range fields {
			if !signature.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != signature.FieldID {
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
	if value, ok := suo.mutation.KeyID(); ok {
		_spec.SetField(signature.FieldKeyID, field.TypeString, value)
	}
	if value, ok := suo.mutation.Signature(); ok {
		_spec.SetField(signature.FieldSignature, field.TypeString, value)
	}
	if suo.mutation.DsseCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   signature.DsseTable,
			Columns: []string{signature.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(dsse.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.DsseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   signature.DsseTable,
			Columns: []string{signature.DsseColumn},
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
	if suo.mutation.TimestampsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   signature.TimestampsTable,
			Columns: []string{signature.TimestampsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.RemovedTimestampsIDs(); len(nodes) > 0 && !suo.mutation.TimestampsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   signature.TimestampsTable,
			Columns: []string{signature.TimestampsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := suo.mutation.TimestampsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   signature.TimestampsTable,
			Columns: []string{signature.TimestampsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(timestamp.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Signature{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{signature.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	suo.mutation.done = true
	return _node, nil
}
