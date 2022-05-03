// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/testifysec/archivist/ent/dsse"
	"github.com/testifysec/archivist/ent/signature"
)

// SignatureCreate is the builder for creating a Signature entity.
type SignatureCreate struct {
	config
	mutation *SignatureMutation
	hooks    []Hook
}

// SetKeyID sets the "key_id" field.
func (sc *SignatureCreate) SetKeyID(s string) *SignatureCreate {
	sc.mutation.SetKeyID(s)
	return sc
}

// SetSignature sets the "signature" field.
func (sc *SignatureCreate) SetSignature(s string) *SignatureCreate {
	sc.mutation.SetSignature(s)
	return sc
}

// SetDsseID sets the "dsse" edge to the Dsse entity by ID.
func (sc *SignatureCreate) SetDsseID(id int) *SignatureCreate {
	sc.mutation.SetDsseID(id)
	return sc
}

// SetNillableDsseID sets the "dsse" edge to the Dsse entity by ID if the given value is not nil.
func (sc *SignatureCreate) SetNillableDsseID(id *int) *SignatureCreate {
	if id != nil {
		sc = sc.SetDsseID(*id)
	}
	return sc
}

// SetDsse sets the "dsse" edge to the Dsse entity.
func (sc *SignatureCreate) SetDsse(d *Dsse) *SignatureCreate {
	return sc.SetDsseID(d.ID)
}

// Mutation returns the SignatureMutation object of the builder.
func (sc *SignatureCreate) Mutation() *SignatureMutation {
	return sc.mutation
}

// Save creates the Signature in the database.
func (sc *SignatureCreate) Save(ctx context.Context) (*Signature, error) {
	var (
		err  error
		node *Signature
	)
	if len(sc.hooks) == 0 {
		if err = sc.check(); err != nil {
			return nil, err
		}
		node, err = sc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SignatureMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = sc.check(); err != nil {
				return nil, err
			}
			sc.mutation = mutation
			if node, err = sc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(sc.hooks) - 1; i >= 0; i-- {
			if sc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = sc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, sc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SignatureCreate) SaveX(ctx context.Context) *Signature {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SignatureCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SignatureCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SignatureCreate) check() error {
	if _, ok := sc.mutation.KeyID(); !ok {
		return &ValidationError{Name: "key_id", err: errors.New(`ent: missing required field "Signature.key_id"`)}
	}
	if v, ok := sc.mutation.KeyID(); ok {
		if err := signature.KeyIDValidator(v); err != nil {
			return &ValidationError{Name: "key_id", err: fmt.Errorf(`ent: validator failed for field "Signature.key_id": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Signature(); !ok {
		return &ValidationError{Name: "signature", err: errors.New(`ent: missing required field "Signature.signature"`)}
	}
	if v, ok := sc.mutation.Signature(); ok {
		if err := signature.SignatureValidator(v); err != nil {
			return &ValidationError{Name: "signature", err: fmt.Errorf(`ent: validator failed for field "Signature.signature": %w`, err)}
		}
	}
	return nil
}

func (sc *SignatureCreate) sqlSave(ctx context.Context) (*Signature, error) {
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (sc *SignatureCreate) createSpec() (*Signature, *sqlgraph.CreateSpec) {
	var (
		_node = &Signature{config: sc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: signature.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: signature.FieldID,
			},
		}
	)
	if value, ok := sc.mutation.KeyID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: signature.FieldKeyID,
		})
		_node.KeyID = value
	}
	if value, ok := sc.mutation.Signature(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: signature.FieldSignature,
		})
		_node.Signature = value
	}
	if nodes := sc.mutation.DsseIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   signature.DsseTable,
			Columns: []string{signature.DsseColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: dsse.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.dsse_signatures = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SignatureCreateBulk is the builder for creating many Signature entities in bulk.
type SignatureCreateBulk struct {
	config
	builders []*SignatureCreate
}

// Save creates the Signature entities in the database.
func (scb *SignatureCreateBulk) Save(ctx context.Context) ([]*Signature, error) {
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Signature, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SignatureMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{err.Error(), err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SignatureCreateBulk) SaveX(ctx context.Context) []*Signature {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SignatureCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SignatureCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}
