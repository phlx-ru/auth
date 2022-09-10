// Code generated by ent, DO NOT EDIT.

package ent

import (
	"auth/ent/code"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// Code is the model entity for the Code schema.
type Code struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// user identification number
	UserID int `json:"user_id,omitempty"`
	// code content, ex.: 1234
	Content string `json:"-"`
	// creation time of code
	CreatedAt time.Time `json:"created_at,omitempty"`
	// last update time of code
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// time of code expiration
	ExpiredAt time.Time `json:"expired_at,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Code) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case code.FieldID, code.FieldUserID:
			values[i] = new(sql.NullInt64)
		case code.FieldContent:
			values[i] = new(sql.NullString)
		case code.FieldCreatedAt, code.FieldUpdatedAt, code.FieldExpiredAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Code", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Code fields.
func (c *Code) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case code.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case code.FieldUserID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field user_id", values[i])
			} else if value.Valid {
				c.UserID = int(value.Int64)
			}
		case code.FieldContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value.Valid {
				c.Content = value.String
			}
		case code.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				c.CreatedAt = value.Time
			}
		case code.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				c.UpdatedAt = value.Time
			}
		case code.FieldExpiredAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field expired_at", values[i])
			} else if value.Valid {
				c.ExpiredAt = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Code.
// Note that you need to call Code.Unwrap() before calling this method if this Code
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Code) Update() *CodeUpdateOne {
	return (&CodeClient{config: c.config}).UpdateOne(c)
}

// Unwrap unwraps the Code entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Code) Unwrap() *Code {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Code is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Code) String() string {
	var builder strings.Builder
	builder.WriteString("Code(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("user_id=")
	builder.WriteString(fmt.Sprintf("%v", c.UserID))
	builder.WriteString(", ")
	builder.WriteString("content=<sensitive>")
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(c.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(c.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("expired_at=")
	builder.WriteString(c.ExpiredAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// Codes is a parsable slice of Code.
type Codes []*Code

func (c Codes) config(cfg config) {
	for _i := range c {
		c[_i].config = cfg
	}
}
