package tweetforward

import (
	"context"
	"errors"
)

var ErrEmpty = errors.New("empty context")

// new tweet context
type NewTweetContext struct {
	Chatid   int64
	Tweeturl string
}

// put NewTweetContext to context.Context
func GetNewTweetFromContext(ctx context.Context) (ret NewTweetContext, err error) {
	v := ctx.Value("newtweet")
	if v == nil {
		err = ErrEmpty
		return
	}

	ret = v.(NewTweetContext)
	return
}

// get NewTweetContext from context.Context
func PutNewTweetToContext(ctx context.Context, v NewTweetContext) context.Context {
	val := context.WithValue(context.Background(), "newtweet", v)

	return val
}
