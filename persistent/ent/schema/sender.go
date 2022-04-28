package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Sender holds the schema definition for the Sender entity.
type Sender struct {
	ent.Schema
}

// Fields of the Sender.
func (Sender) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("username"),
	}
}

// Edges of the Sender.
func (Sender) Edges() []ent.Edge {
	return nil
}
