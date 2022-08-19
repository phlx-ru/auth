package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String(`display_name`).
			Comment(`user displayed name`),

		field.String(`type`).
			Comment(`user category (admin|dispatcher|driver)`),

		field.String(`email`).
			Optional().
			Nillable().
			Comment(`users email`),

		field.String(`phone`).
			Optional().
			Nillable().
			Comment(`users phone`),

		field.String(`password_hash`).
			Optional().
			Nillable().
			Sensitive().
			Comment(`hashed string with password`),

		field.String(`password_reset`).
			Optional().
			Nillable().
			Sensitive().
			Comment(`reset string for change password`),

		field.Time(`created_at`).
			Default(time.Now).
			Immutable().
			Comment(`creation time of code`),

		field.Time(`updated_at`).
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				&entsql.Annotation{
					Default: `CURRENT_TIMESTAMP`,
				},
			).
			Comment(`last update time of code`),

		field.Time(`deactivated_at`).
			Optional().
			Nillable().
			Default(nil).
			Comment(`user deactivation time`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
