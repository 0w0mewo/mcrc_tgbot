package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// TweetUser holds the schema definition for the TweetUser entity.
type TweetUser struct {
	ent.Schema
}

// Fields of the TweetUser.
func (TweetUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("username"),
	}
}

// Edges of the TweetUser.
func (TweetUser) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.To("ChatTweetSubscription", ChatTweetSubscription.Type),
	}
}

func (TweetUser) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
	}
}
