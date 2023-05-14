// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/tenant"
	"gitlab.com/testifysec/judge-platform/judge-api/ent/user"
)

// TenantCreate is the builder for creating a Tenant entity.
type TenantCreate struct {
	config
	mutation *TenantMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (tc *TenantCreate) SetCreatedAt(t time.Time) *TenantCreate {
	tc.mutation.SetCreatedAt(t)
	return tc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tc *TenantCreate) SetNillableCreatedAt(t *time.Time) *TenantCreate {
	if t != nil {
		tc.SetCreatedAt(*t)
	}
	return tc
}

// SetUpdatedAt sets the "updated_at" field.
func (tc *TenantCreate) SetUpdatedAt(t time.Time) *TenantCreate {
	tc.mutation.SetUpdatedAt(t)
	return tc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tc *TenantCreate) SetNillableUpdatedAt(t *time.Time) *TenantCreate {
	if t != nil {
		tc.SetUpdatedAt(*t)
	}
	return tc
}

// SetName sets the "name" field.
func (tc *TenantCreate) SetName(s string) *TenantCreate {
	tc.mutation.SetName(s)
	return tc
}

// SetDescription sets the "description" field.
func (tc *TenantCreate) SetDescription(s string) *TenantCreate {
	tc.mutation.SetDescription(s)
	return tc
}

// SetType sets the "type" field.
func (tc *TenantCreate) SetType(t tenant.Type) *TenantCreate {
	tc.mutation.SetType(t)
	return tc
}

// SetID sets the "id" field.
func (tc *TenantCreate) SetID(u uuid.UUID) *TenantCreate {
	tc.mutation.SetID(u)
	return tc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tc *TenantCreate) SetNillableID(u *uuid.UUID) *TenantCreate {
	if u != nil {
		tc.SetID(*u)
	}
	return tc
}

// SetCreatedByID sets the "created_by" edge to the User entity by ID.
func (tc *TenantCreate) SetCreatedByID(id uuid.UUID) *TenantCreate {
	tc.mutation.SetCreatedByID(id)
	return tc
}

// SetNillableCreatedByID sets the "created_by" edge to the User entity by ID if the given value is not nil.
func (tc *TenantCreate) SetNillableCreatedByID(id *uuid.UUID) *TenantCreate {
	if id != nil {
		tc = tc.SetCreatedByID(*id)
	}
	return tc
}

// SetCreatedBy sets the "created_by" edge to the User entity.
func (tc *TenantCreate) SetCreatedBy(u *User) *TenantCreate {
	return tc.SetCreatedByID(u.ID)
}

// SetModifiedByID sets the "modified_by" edge to the User entity by ID.
func (tc *TenantCreate) SetModifiedByID(id uuid.UUID) *TenantCreate {
	tc.mutation.SetModifiedByID(id)
	return tc
}

// SetNillableModifiedByID sets the "modified_by" edge to the User entity by ID if the given value is not nil.
func (tc *TenantCreate) SetNillableModifiedByID(id *uuid.UUID) *TenantCreate {
	if id != nil {
		tc = tc.SetModifiedByID(*id)
	}
	return tc
}

// SetModifiedBy sets the "modified_by" edge to the User entity.
func (tc *TenantCreate) SetModifiedBy(u *User) *TenantCreate {
	return tc.SetModifiedByID(u.ID)
}

// AddUserIDs adds the "users" edge to the User entity by IDs.
func (tc *TenantCreate) AddUserIDs(ids ...uuid.UUID) *TenantCreate {
	tc.mutation.AddUserIDs(ids...)
	return tc
}

// AddUsers adds the "users" edges to the User entity.
func (tc *TenantCreate) AddUsers(u ...*User) *TenantCreate {
	ids := make([]uuid.UUID, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return tc.AddUserIDs(ids...)
}

// SetParentID sets the "parent" edge to the Tenant entity by ID.
func (tc *TenantCreate) SetParentID(id uuid.UUID) *TenantCreate {
	tc.mutation.SetParentID(id)
	return tc
}

// SetNillableParentID sets the "parent" edge to the Tenant entity by ID if the given value is not nil.
func (tc *TenantCreate) SetNillableParentID(id *uuid.UUID) *TenantCreate {
	if id != nil {
		tc = tc.SetParentID(*id)
	}
	return tc
}

// SetParent sets the "parent" edge to the Tenant entity.
func (tc *TenantCreate) SetParent(t *Tenant) *TenantCreate {
	return tc.SetParentID(t.ID)
}

// AddChildIDs adds the "children" edge to the Tenant entity by IDs.
func (tc *TenantCreate) AddChildIDs(ids ...uuid.UUID) *TenantCreate {
	tc.mutation.AddChildIDs(ids...)
	return tc
}

// AddChildren adds the "children" edges to the Tenant entity.
func (tc *TenantCreate) AddChildren(t ...*Tenant) *TenantCreate {
	ids := make([]uuid.UUID, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tc.AddChildIDs(ids...)
}

// Mutation returns the TenantMutation object of the builder.
func (tc *TenantCreate) Mutation() *TenantMutation {
	return tc.mutation
}

// Save creates the Tenant in the database.
func (tc *TenantCreate) Save(ctx context.Context) (*Tenant, error) {
	if err := tc.defaults(); err != nil {
		return nil, err
	}
	return withHooks[*Tenant, TenantMutation](ctx, tc.sqlSave, tc.mutation, tc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (tc *TenantCreate) SaveX(ctx context.Context) *Tenant {
	v, err := tc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tc *TenantCreate) Exec(ctx context.Context) error {
	_, err := tc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tc *TenantCreate) ExecX(ctx context.Context) {
	if err := tc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tc *TenantCreate) defaults() error {
	if _, ok := tc.mutation.CreatedAt(); !ok {
		if tenant.DefaultCreatedAt == nil {
			return fmt.Errorf("ent: uninitialized tenant.DefaultCreatedAt (forgotten import ent/runtime?)")
		}
		v := tenant.DefaultCreatedAt()
		tc.mutation.SetCreatedAt(v)
	}
	if _, ok := tc.mutation.UpdatedAt(); !ok {
		if tenant.DefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized tenant.DefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := tenant.DefaultUpdatedAt()
		tc.mutation.SetUpdatedAt(v)
	}
	if _, ok := tc.mutation.ID(); !ok {
		if tenant.DefaultID == nil {
			return fmt.Errorf("ent: uninitialized tenant.DefaultID (forgotten import ent/runtime?)")
		}
		v := tenant.DefaultID()
		tc.mutation.SetID(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (tc *TenantCreate) check() error {
	if _, ok := tc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Tenant.created_at"`)}
	}
	if _, ok := tc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Tenant.updated_at"`)}
	}
	if _, ok := tc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Tenant.name"`)}
	}
	if _, ok := tc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Tenant.description"`)}
	}
	if _, ok := tc.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "Tenant.type"`)}
	}
	if v, ok := tc.mutation.GetType(); ok {
		if err := tenant.TypeValidator(v); err != nil {
			return &ValidationError{Name: "type", err: fmt.Errorf(`ent: validator failed for field "Tenant.type": %w`, err)}
		}
	}
	return nil
}

func (tc *TenantCreate) sqlSave(ctx context.Context) (*Tenant, error) {
	if err := tc.check(); err != nil {
		return nil, err
	}
	_node, _spec := tc.createSpec()
	if err := sqlgraph.CreateNode(ctx, tc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	tc.mutation.id = &_node.ID
	tc.mutation.done = true
	return _node, nil
}

func (tc *TenantCreate) createSpec() (*Tenant, *sqlgraph.CreateSpec) {
	var (
		_node = &Tenant{config: tc.config}
		_spec = sqlgraph.NewCreateSpec(tenant.Table, sqlgraph.NewFieldSpec(tenant.FieldID, field.TypeUUID))
	)
	if id, ok := tc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tc.mutation.CreatedAt(); ok {
		_spec.SetField(tenant.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := tc.mutation.UpdatedAt(); ok {
		_spec.SetField(tenant.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := tc.mutation.Name(); ok {
		_spec.SetField(tenant.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := tc.mutation.Description(); ok {
		_spec.SetField(tenant.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if value, ok := tc.mutation.GetType(); ok {
		_spec.SetField(tenant.FieldType, field.TypeEnum, value)
		_node.Type = value
	}
	if nodes := tc.mutation.CreatedByIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   tenant.CreatedByTable,
			Columns: []string{tenant.CreatedByColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.tenant_created_by = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.ModifiedByIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: false,
			Table:   tenant.ModifiedByTable,
			Columns: []string{tenant.ModifiedByColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.tenant_modified_by = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.UsersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   tenant.UsersTable,
			Columns: tenant.UsersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   tenant.ParentTable,
			Columns: []string{tenant.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tenant.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.tenant_children = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tc.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   tenant.ChildrenTable,
			Columns: []string{tenant.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tenant.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TenantCreateBulk is the builder for creating many Tenant entities in bulk.
type TenantCreateBulk struct {
	config
	builders []*TenantCreate
}

// Save creates the Tenant entities in the database.
func (tcb *TenantCreateBulk) Save(ctx context.Context) ([]*Tenant, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tcb.builders))
	nodes := make([]*Tenant, len(tcb.builders))
	mutators := make([]Mutator, len(tcb.builders))
	for i := range tcb.builders {
		func(i int, root context.Context) {
			builder := tcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TenantMutation)
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
					_, err = mutators[i+1].Mutate(root, tcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tcb *TenantCreateBulk) SaveX(ctx context.Context) []*Tenant {
	v, err := tcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tcb *TenantCreateBulk) Exec(ctx context.Context) error {
	_, err := tcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tcb *TenantCreateBulk) ExecX(ctx context.Context) {
	if err := tcb.Exec(ctx); err != nil {
		panic(err)
	}
}
