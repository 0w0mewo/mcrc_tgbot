package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ChatTweetSubscription holds the schema definition for the ChatTweetSubscription entity.
type ChatTweetSubscription struct {
	ent.Schema
}

// Fields of the ChatTweetSubscription.
func (ChatTweetSubscription) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("chat_id").Optional(),
		field.String("tweeter_id").Optional(),
		field.String("last_tweet"),
	}
}

// Edges of the ChatTweetSubscription.
func (ChatTweetSubscription) Edges() []ent.Edge {
	return []ent.Edge{
		// edge.From("TweetUser", TweetUser.Type).Ref("ChatTweetSubscription"),
		edge.To("subscribed_tweeter", TweetUser.Type).Field("tweeter_id").Unique(),
		edge.To("subscribed_chat", Chat.Type).Field("chat_id").Unique(),
	}
}

func (ChatTweetSubscription) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("chat_id", "tweeter_id"),
	}
}
