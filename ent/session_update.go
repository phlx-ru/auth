// Code generated by ent, DO NOT EDIT.

package ent

import (
	"auth/ent/predicate"
	"auth/ent/session"
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SessionUpdate is the builder for updating Session entities.
type SessionUpdate struct {
	config
	hooks    []Hook
	mutation *SessionMutation
}

// Where appends a list predicates to the SessionUpdate builder.
func (su *SessionUpdate) Where(ps ...predicate.Session) *SessionUpdate {
	su.mutation.Where(ps...)
	return su
}

// SetUserID sets the "user_id" field.
func (su *SessionUpdate) SetUserID(i int) *SessionUpdate {
	su.mutation.ResetUserID()
	su.mutation.SetUserID(i)
	return su
}

// AddUserID adds i to the "user_id" field.
func (su *SessionUpdate) AddUserID(i int) *SessionUpdate {
	su.mutation.AddUserID(i)
	return su
}

// SetToken sets the "token" field.
func (su *SessionUpdate) SetToken(s string) *SessionUpdate {
	su.mutation.SetToken(s)
	return su
}

// SetIP sets the "ip" field.
func (su *SessionUpdate) SetIP(s string) *SessionUpdate {
	su.mutation.SetIP(s)
	return su
}

// SetUserAgent sets the "user_agent" field.
func (su *SessionUpdate) SetUserAgent(s string) *SessionUpdate {
	su.mutation.SetUserAgent(s)
	return su
}

// SetDeviceID sets the "device_id" field.
func (su *SessionUpdate) SetDeviceID(s string) *SessionUpdate {
	su.mutation.SetDeviceID(s)
	return su
}

// SetNillableDeviceID sets the "device_id" field if the given value is not nil.
func (su *SessionUpdate) SetNillableDeviceID(s *string) *SessionUpdate {
	if s != nil {
		su.SetDeviceID(*s)
	}
	return su
}

// ClearDeviceID clears the value of the "device_id" field.
func (su *SessionUpdate) ClearDeviceID() *SessionUpdate {
	su.mutation.ClearDeviceID()
	return su
}

// SetUpdatedAt sets the "updated_at" field.
func (su *SessionUpdate) SetUpdatedAt(t time.Time) *SessionUpdate {
	su.mutation.SetUpdatedAt(t)
	return su
}

// SetExpiredAt sets the "expired_at" field.
func (su *SessionUpdate) SetExpiredAt(t time.Time) *SessionUpdate {
	su.mutation.SetExpiredAt(t)
	return su
}

// SetNillableExpiredAt sets the "expired_at" field if the given value is not nil.
func (su *SessionUpdate) SetNillableExpiredAt(t *time.Time) *SessionUpdate {
	if t != nil {
		su.SetExpiredAt(*t)
	}
	return su
}

// SetIsActive sets the "is_active" field.
func (su *SessionUpdate) SetIsActive(b bool) *SessionUpdate {
	su.mutation.SetIsActive(b)
	return su
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (su *SessionUpdate) SetNillableIsActive(b *bool) *SessionUpdate {
	if b != nil {
		su.SetIsActive(*b)
	}
	return su
}

// Mutation returns the SessionMutation object of the builder.
func (su *SessionUpdate) Mutation() *SessionMutation {
	return su.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (su *SessionUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	su.defaults()
	if len(su.hooks) == 0 {
		affected, err = su.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SessionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			su.mutation = mutation
			affected, err = su.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(su.hooks) - 1; i >= 0; i-- {
			if su.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = su.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, su.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (su *SessionUpdate) SaveX(ctx context.Context) int {
	affected, err := su.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (su *SessionUpdate) Exec(ctx context.Context) error {
	_, err := su.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (su *SessionUpdate) ExecX(ctx context.Context) {
	if err := su.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (su *SessionUpdate) defaults() {
	if _, ok := su.mutation.UpdatedAt(); !ok {
		v := session.UpdateDefaultUpdatedAt()
		su.mutation.SetUpdatedAt(v)
	}
}

func (su *SessionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   session.Table,
			Columns: session.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: session.FieldID,
			},
		},
	}
	if ps := su.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := su.mutation.UserID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: session.FieldUserID,
		})
	}
	if value, ok := su.mutation.AddedUserID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: session.FieldUserID,
		})
	}
	if value, ok := su.mutation.Token(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldToken,
		})
	}
	if value, ok := su.mutation.IP(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldIP,
		})
	}
	if value, ok := su.mutation.UserAgent(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldUserAgent,
		})
	}
	if value, ok := su.mutation.DeviceID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldDeviceID,
		})
	}
	if su.mutation.DeviceIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: session.FieldDeviceID,
		})
	}
	if value, ok := su.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: session.FieldUpdatedAt,
		})
	}
	if value, ok := su.mutation.ExpiredAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: session.FieldExpiredAt,
		})
	}
	if value, ok := su.mutation.IsActive(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: session.FieldIsActive,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, su.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{session.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	return n, nil
}

// SessionUpdateOne is the builder for updating a single Session entity.
type SessionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SessionMutation
}

// SetUserID sets the "user_id" field.
func (suo *SessionUpdateOne) SetUserID(i int) *SessionUpdateOne {
	suo.mutation.ResetUserID()
	suo.mutation.SetUserID(i)
	return suo
}

// AddUserID adds i to the "user_id" field.
func (suo *SessionUpdateOne) AddUserID(i int) *SessionUpdateOne {
	suo.mutation.AddUserID(i)
	return suo
}

// SetToken sets the "token" field.
func (suo *SessionUpdateOne) SetToken(s string) *SessionUpdateOne {
	suo.mutation.SetToken(s)
	return suo
}

// SetIP sets the "ip" field.
func (suo *SessionUpdateOne) SetIP(s string) *SessionUpdateOne {
	suo.mutation.SetIP(s)
	return suo
}

// SetUserAgent sets the "user_agent" field.
func (suo *SessionUpdateOne) SetUserAgent(s string) *SessionUpdateOne {
	suo.mutation.SetUserAgent(s)
	return suo
}

// SetDeviceID sets the "device_id" field.
func (suo *SessionUpdateOne) SetDeviceID(s string) *SessionUpdateOne {
	suo.mutation.SetDeviceID(s)
	return suo
}

// SetNillableDeviceID sets the "device_id" field if the given value is not nil.
func (suo *SessionUpdateOne) SetNillableDeviceID(s *string) *SessionUpdateOne {
	if s != nil {
		suo.SetDeviceID(*s)
	}
	return suo
}

// ClearDeviceID clears the value of the "device_id" field.
func (suo *SessionUpdateOne) ClearDeviceID() *SessionUpdateOne {
	suo.mutation.ClearDeviceID()
	return suo
}

// SetUpdatedAt sets the "updated_at" field.
func (suo *SessionUpdateOne) SetUpdatedAt(t time.Time) *SessionUpdateOne {
	suo.mutation.SetUpdatedAt(t)
	return suo
}

// SetExpiredAt sets the "expired_at" field.
func (suo *SessionUpdateOne) SetExpiredAt(t time.Time) *SessionUpdateOne {
	suo.mutation.SetExpiredAt(t)
	return suo
}

// SetNillableExpiredAt sets the "expired_at" field if the given value is not nil.
func (suo *SessionUpdateOne) SetNillableExpiredAt(t *time.Time) *SessionUpdateOne {
	if t != nil {
		suo.SetExpiredAt(*t)
	}
	return suo
}

// SetIsActive sets the "is_active" field.
func (suo *SessionUpdateOne) SetIsActive(b bool) *SessionUpdateOne {
	suo.mutation.SetIsActive(b)
	return suo
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (suo *SessionUpdateOne) SetNillableIsActive(b *bool) *SessionUpdateOne {
	if b != nil {
		suo.SetIsActive(*b)
	}
	return suo
}

// Mutation returns the SessionMutation object of the builder.
func (suo *SessionUpdateOne) Mutation() *SessionMutation {
	return suo.mutation
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (suo *SessionUpdateOne) Select(field string, fields ...string) *SessionUpdateOne {
	suo.fields = append([]string{field}, fields...)
	return suo
}

// Save executes the query and returns the updated Session entity.
func (suo *SessionUpdateOne) Save(ctx context.Context) (*Session, error) {
	var (
		err  error
		node *Session
	)
	suo.defaults()
	if len(suo.hooks) == 0 {
		node, err = suo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*SessionMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			suo.mutation = mutation
			node, err = suo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(suo.hooks) - 1; i >= 0; i-- {
			if suo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = suo.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, suo.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*Session)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from SessionMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (suo *SessionUpdateOne) SaveX(ctx context.Context) *Session {
	node, err := suo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (suo *SessionUpdateOne) Exec(ctx context.Context) error {
	_, err := suo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (suo *SessionUpdateOne) ExecX(ctx context.Context) {
	if err := suo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (suo *SessionUpdateOne) defaults() {
	if _, ok := suo.mutation.UpdatedAt(); !ok {
		v := session.UpdateDefaultUpdatedAt()
		suo.mutation.SetUpdatedAt(v)
	}
}

func (suo *SessionUpdateOne) sqlSave(ctx context.Context) (_node *Session, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   session.Table,
			Columns: session.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: session.FieldID,
			},
		},
	}
	id, ok := suo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Session.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := suo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, session.FieldID)
		for _, f := range fields {
			if !session.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != session.FieldID {
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
	if value, ok := suo.mutation.UserID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: session.FieldUserID,
		})
	}
	if value, ok := suo.mutation.AddedUserID(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: session.FieldUserID,
		})
	}
	if value, ok := suo.mutation.Token(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldToken,
		})
	}
	if value, ok := suo.mutation.IP(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldIP,
		})
	}
	if value, ok := suo.mutation.UserAgent(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldUserAgent,
		})
	}
	if value, ok := suo.mutation.DeviceID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: session.FieldDeviceID,
		})
	}
	if suo.mutation.DeviceIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: session.FieldDeviceID,
		})
	}
	if value, ok := suo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: session.FieldUpdatedAt,
		})
	}
	if value, ok := suo.mutation.ExpiredAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: session.FieldExpiredAt,
		})
	}
	if value, ok := suo.mutation.IsActive(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeBool,
			Value:  value,
			Column: session.FieldIsActive,
		})
	}
	_node = &Session{config: suo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, suo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{session.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	return _node, nil
}
