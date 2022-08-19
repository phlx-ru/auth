package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.Int(`user_id`).
			Comment(`user identification number`),

		field.String(`token`).
			Sensitive().
			Comment(`user auth token`),

		field.Time(`created_at`).
			Default(time.Now).
			Immutable().
			Comment(`creation time of session`),

		field.Time(`updated_at`).
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				&entsql.Annotation{
					Default: `CURRENT_TIMESTAMP`,
				},
			).
			Comment(`last update time of session`),

		field.Time(`expired_at`).
			Default(time.Now).
			Comment(`time of session expiration`),

		field.Bool(`is_active`).
			Default(true).
			Comment(`session activity flag`),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return nil
}
