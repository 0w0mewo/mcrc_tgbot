package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Chat holds the schema definition for the Chat entity.
type Chat struct {
	ent.Schema
}

// Fields of the Chat.
func (Chat) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("name"),
	}
}

// Edges of the Chat.
func (Chat) Edges() []ent.Edge {
	return nil
}
