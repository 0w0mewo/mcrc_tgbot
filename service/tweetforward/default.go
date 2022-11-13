package tweetforward

import (
	models "github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/0w0mewo/mcrc_tgbot/utils"
)

var defaultForwarder *tweetForwarder

func init() {
	defaultForwarder = newTweetForwarder()

}

func Subscribe(chat models.Chat, username string) error {
	return defaultForwarder.Subscribe(chat, username)
}

func Unsubscribe(chatid int64, username string) error {
	return defaultForwarder.Unsubscribe(chatid, username)
}

func GetSubscriptions(chatid int64) (models.TweetUserSlice, error) {
	return defaultForwarder.GetSubscriptions(chatid)
}

func DoWhenSet(cb utils.EventCallback) {
	defaultForwarder.DoWhenSet(cb)
}

func Shutdown() {
	defaultForwarder.Shutdown()
}
