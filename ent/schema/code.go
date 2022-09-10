package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Code holds the schema definition for the Code entity.
type Code struct {
	ent.Schema
}

// Fields of the Code.
func (Code) Fields() []ent.Field {
	return []ent.Field{
		field.Int(`user_id`).
			Comment(`user identification number`),

		field.String(`content`).
			Sensitive().
			Comment(`code content, ex.: 1234`),

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

		field.Time(`expired_at`).
			Default(time.Now).
			Comment(`time of code expiration`),
	}
}

// Edges of the Code.
func (Code) Edges() []ent.Edge {
	return nil
}

func (Code) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(`user_id`),
		index.Fields(`expired_at`),
	}
}
