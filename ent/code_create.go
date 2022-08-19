// Code generated by ent, DO NOT EDIT.

package ent

import (
	"auth/ent/code"
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CodeCreate is the builder for creating a Code entity.
type CodeCreate struct {
	config
	mutation *CodeMutation
	hooks    []Hook
}

// SetUserID sets the "user_id" field.
func (cc *CodeCreate) SetUserID(i int) *CodeCreate {
	cc.mutation.SetUserID(i)
	return cc
}

// SetContent sets the "content" field.
func (cc *CodeCreate) SetContent(s string) *CodeCreate {
	cc.mutation.SetContent(s)
	return cc
}

// SetCreatedAt sets the "created_at" field.
func (cc *CodeCreate) SetCreatedAt(t time.Time) *CodeCreate {
	cc.mutation.SetCreatedAt(t)
	return cc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (cc *CodeCreate) SetNillableCreatedAt(t *time.Time) *CodeCreate {
	if t != nil {
		cc.SetCreatedAt(*t)
	}
	return cc
}

// SetUpdatedAt sets the "updated_at" field.
func (cc *CodeCreate) SetUpdatedAt(t time.Time) *CodeCreate {
	cc.mutation.SetUpdatedAt(t)
	return cc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (cc *CodeCreate) SetNillableUpdatedAt(t *time.Time) *CodeCreate {
	if t != nil {
		cc.SetUpdatedAt(*t)
	}
	return cc
}

// SetExpiredAt sets the "expired_at" field.
func (cc *CodeCreate) SetExpiredAt(t time.Time) *CodeCreate {
	cc.mutation.SetExpiredAt(t)
	return cc
}

// SetNillableExpiredAt sets the "expired_at" field if the given value is not nil.
func (cc *CodeCreate) SetNillableExpiredAt(t *time.Time) *CodeCreate {
	if t != nil {
		cc.SetExpiredAt(*t)
	}
	return cc
}

// SetRetries sets the "retries" field.
func (cc *CodeCreate) SetRetries(i int) *CodeCreate {
	cc.mutation.SetRetries(i)
	return cc
}

// SetNillableRetries sets the "retries" field if the given value is not nil.
func (cc *CodeCreate) SetNillableRetries(i *int) *CodeCreate {
	if i != nil {
		cc.SetRetries(*i)
	}
	return cc
}

// Mutation returns the CodeMutation object of the builder.
func (cc *CodeCreate) Mutation() *CodeMutation {
	return cc.mutation
}

// Save creates the Code in the database.
func (cc *CodeCreate) Save(ctx context.Context) (*Code, error) {
	var (
		err  error
		node *Code
	)
	cc.defaults()
	if len(cc.hooks) == 0 {
		if err = cc.check(); err != nil {
			return nil, err
		}
		node, err = cc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*CodeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = cc.check(); err != nil {
				return nil, err
			}
			cc.mutation = mutation
			if node, err = cc.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(cc.hooks) - 1; i >= 0; i-- {
			if cc.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = cc.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, cc.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Code)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from CodeMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (cc *CodeCreate) SaveX(ctx context.Context) *Code {
	v, err := cc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (cc *CodeCreate) Exec(ctx context.Context) error {
	_, err := cc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cc *CodeCreate) ExecX(ctx context.Context) {
	if err := cc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (cc *CodeCreate) defaults() {
	if _, ok := cc.mutation.CreatedAt(); !ok {
		v := code.DefaultCreatedAt()
		cc.mutation.SetCreatedAt(v)
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		v := code.DefaultUpdatedAt()
		cc.mutation.SetUpdatedAt(v)
	}
	if _, ok := cc.mutation.ExpiredAt(); !ok {
		v := code.DefaultExpiredAt()
		cc.mutation.SetExpiredAt(v)
	}
	if _, ok := cc.mutation.Retries(); !ok {
		v := code.DefaultRetries
		cc.mutation.SetRetries(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (cc *CodeCreate) check() error {
	if _, ok := cc.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "Code.user_id"`)}
	}
	if _, ok := cc.mutation.Content(); !ok {
		return &ValidationError{Name: "content", err: errors.New(`ent: missing required field "Code.content"`)}
	}
	if _, ok := cc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "Code.created_at"`)}
	}
	if _, ok := cc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "Code.updated_at"`)}
	}
	if _, ok := cc.mutation.ExpiredAt(); !ok {
		return &ValidationError{Name: "expired_at", err: errors.New(`ent: missing required field "Code.expired_at"`)}
	}
	if _, ok := cc.mutation.Retries(); !ok {
		return &ValidationError{Name: "retries", err: errors.New(`ent: missing required field "Code.retries"`)}
	}
	return nil
}

func (cc *CodeCreate) sqlSave(ctx context.Context) (*Code, error) {
	_node, _spec := cc.createSpec()
	if err := sqlgraph.CreateNode(ctx, cc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (cc *CodeCreate) createSpec() (*Code, *sqlgraph.CreateSpec) {
	var (
		_node = &Code{config: cc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: code.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: code.FieldID,
			},
		}
	)
	if value, ok := cc.mutation.UserID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: code.FieldUserID,
		})
		_node.UserID = value
	}
	if value, ok := cc.mutation.Content(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: code.FieldContent,
		})
		_node.Content = value
	}
	if value, ok := cc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: code.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := cc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: code.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := cc.mutation.ExpiredAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: code.FieldExpiredAt,
		})
		_node.ExpiredAt = value
	}
	if value, ok := cc.mutation.Retries(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: code.FieldRetries,
		})
		_node.Retries = value
	}
	return _node, _spec
}

// CodeCreateBulk is the builder for creating many Code entities in bulk.
type CodeCreateBulk struct {
	config
	builders []*CodeCreate
}

// Save creates the Code entities in the database.
func (ccb *CodeCreateBulk) Save(ctx context.Context) ([]*Code, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ccb.builders))
	nodes := make([]*Code, len(ccb.builders))
	mutators := make([]Mutator, len(ccb.builders))
	for i := range ccb.builders {
		func(i int, root context.Context) {
			builder := ccb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*CodeMutation)
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
					_, err = mutators[i+1].Mutate(root, ccb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ccb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
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
		if _, err := mutators[0].Mutate(ctx, ccb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (ccb *CodeCreateBulk) SaveX(ctx context.Context) []*Code {
	v, err := ccb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (ccb *CodeCreateBulk) Exec(ctx context.Context) error {
	_, err := ccb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (ccb *CodeCreateBulk) ExecX(ctx context.Context) {
	if err := ccb.Exec(ctx); err != nil {
		panic(err)
	}
}
