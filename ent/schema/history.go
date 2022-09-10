package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

const (
	EventLoginOk              = `login.ok`
	EventLoginFailed          = `login.failed`
	EventGenerateCodeRequest  = `generate_code.request`
	EventLoginByCodeOk        = `login_by_code.ok`
	EventLoginByCodeFailed    = `login_by_code.failed`
	EventResetPasswordRequest = `reset_password.request`
	EventNewPasswordOk        = `new_password.ok`
	EventNewPasswordFailed    = `new_password.failed`
	EventChangePasswordOk     = `change_password.ok`
	EventChangePasswordFailed = `change_password.failed`
)

var (
	EventTypes = []string{
		EventLoginOk,
		EventLoginFailed,
		EventGenerateCodeRequest,
		EventLoginByCodeOk,
		EventLoginByCodeFailed,
		EventResetPasswordRequest,
		EventNewPasswordOk,
		EventNewPasswordFailed,
		EventChangePasswordOk,
		EventChangePasswordFailed,
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
			Comment(`type of event`),

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

func (History) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`user_id`),
		index.Fields(`event`),
		index.Fields(`created_at`),
	}
}
