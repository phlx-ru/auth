// Code generated by ent, DO NOT EDIT.

package ent

import (
	"auth/ent/code"
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CodeCreate is the builder for creating a Code entity.
type CodeCreate struct {
	config
	mutation *CodeMutation
	hooks    []Hook
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

// check runs all checks and user-defined validators on the builder.
func (cc *CodeCreate) check() error {
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
