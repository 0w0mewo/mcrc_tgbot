package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.Bytes("msg"),
		field.Int("type"),
		field.Time("timestamp"),
		field.Int64("chat_id").Nillable().Optional(),
		field.Int64("sender_id").Nillable().Optional(),
	}
}

func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		// one sender and one chat (union) belongs to one message
		edge.To("from_chat", Chat.Type).Field("chat_id").Unique(),
		edge.To("from_sender", Sender.Type).Field("sender_id").Unique(),
	}
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("chat_id", "sender_id"),
		index.Fields("chat_id"),
	}
}
