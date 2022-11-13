package tweetforward

import (
	"context"
	"errors"
)

var ErrEmpty = errors.New("empty context")

type NewTweetContext struct {
	Chatid   int64
	Tweeturl string
}

func GetNewTweetFromContext(ctx context.Context) (ret NewTweetContext, err error) {
	v := ctx.Value("newtweet")
	if v == nil {
		err = ErrEmpty
		return
	}

	ret = v.(NewTweetContext)
	return
}

func PutNewTweetToContext(ctx context.Context, v NewTweetContext) context.Context {
	val := context.WithValue(context.Background(), "newtweet", v)

	return val
}
