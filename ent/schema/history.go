package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

const (
	EventLoginOk              = `login.ok`
	EventLoginFailed          = `login.failed`
	EventRequestCode          = `login_by_code.request`
	EventLoginByCodeOk        = `login_by_code.ok`
	EventLoginByCodeFailed    = `login_by_code.failed`
	EventResetPasswordRequest = `reset_password.request`
	EventResetPasswordOk      = `reset_password.ok`
	EventResetPasswordFailed  = `reset_password.failed`
)

var (
	EventTypes = []string{
		EventLoginOk,
		EventLoginFailed,
		EventRequestCode,
		EventLoginByCodeOk,
		EventLoginByCodeFailed,
		EventResetPasswordRequest,
		EventResetPasswordOk,
		EventResetPasswordFailed,
	}
)

// History holds the schema definition for the History entity.
type History struct {
	ent.Schema
}

// Fields of the History.
func (History) Fields() []ent.Field {
	return []ent.Field{
		field.Int(`user_id`).
			Comment(`user identification number`),

		field.Time(`created_at`).
			Default(time.Now).
			Immutable().
			Comment(`creation time of code`),

		field.String(`event`).
			Immutable().
			Comment(`type of event `),

		field.String(`ip`).
			Optional().
			Nillable().
			Comment(`ip of request`),

		field.String(`user_agent`).
			Optional().
			Nillable().
			Comment(`user-agent of request`),
	}
}

// Edges of the History.
func (History) Edges() []ent.Edge {
	return nil
}
