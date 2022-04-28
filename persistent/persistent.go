package persistent

import (
	"context"
	"errors"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/model"
)

var ErrExist = errors.New("already existed")
var ErrNotExist = errors.New("not existed")

// repo for storing telegram messages
type StoredTeleMsgRepo interface {
	StoreMsg(ctx context.Context, chatid int64, chatname string,
		senderid int64, sendername string,
		msg []byte, msgType int, timestamp time.Time) error // store one message
	GetChatMessages(ctx context.Context, chatid int64, page int, pagesize int) ([]*model.TeleMsg, error) // get range of messages
	Count(ctx context.Context) (int, error)                                                              //
	Close() error                                                                                        // close database connection
}

// repo for accessing info on subscribed tweeters
type ChatTweetSubRepo interface {
	Create(ctx context.Context, fromchat model.Chat, sub model.TweetUser) error              // subscribe a tweeter in a chat
	Remove(ctx context.Context, chatid int64, twuid string) error                            // remove subscrition
	Exist(ctx context.Context, chatid int64, twuid string) (bool, error)                     // is subscribed
	GetAllChatSub(ctx context.Context, chatid int64) ([]*model.ChatTweetSubscription, error) // get all subscritions of a chat
	GetAllChat(ctx context.Context) ([]int64, error)                                         // get all chats which subscribe some tweeters
	GetAllSubscribeeByChatId(ctx context.Context, chatid int64) ([]*model.TweetUser, error)
	GetLastTweet(ctx context.Context, chatid int64, twuid string) (string, error)           // get the latest tweet url
	UpdateLastTweet(ctx context.Context, chatid int64, twuid string, newtweet string) error // update the latest tweet url
	GetTweeterOfChatSub(ctx context.Context, id int) (*model.TweetUser, error)              // get tweeter by chatsub table id
	Close() error                                                                           // close database connection
}
