package tweetforward

import "github.com/0w0mewo/mcrc_tgbot/model"

var defaultForwarder *tweetForwarder

func init() {
	defaultForwarder = newTweetForwarder()

}

func Subscribe(chat model.Chat, username string) error {
	return defaultForwarder.Subscribe(chat, username)
}

func Unsubscribe(chatid int64, username string) error {
	return defaultForwarder.Unsubscribe(chatid, username)
}

func GetSubscriptions(chatid int64) ([]*model.TweetUser, error) {
	return defaultForwarder.GetSubscriptions(chatid)
}

func DoWhenSet(cb any) {
	defaultForwarder.DoWhenSet(cb)
}

func Shutdown() {
	defaultForwarder.Shutdown()
}
