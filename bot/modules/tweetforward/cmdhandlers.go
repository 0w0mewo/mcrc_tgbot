package notifier

import (
	"fmt"
	"strings"

	models "github.com/0w0mewo/mcrc_tgbot/model"
	tweetforward "github.com/0w0mewo/mcrc_tgbot/service/tweetforward"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

var ErrNotEnoughArguments = errors.New("not enough args")

func (t *Notifier) tweetSub(c telebot.Context) error {

	var twuname string

	args := c.Args()
	switch {
	case len(args) == 1:
		if args[0] == "--list" || args[0] == "-l" {
			return t.tweetList(c)
		} else {
			twuname = args[0]
		}

	default:
		c.Send("usage: /tweetsub <twitter username>")
		return nil
	}

	var name string
	chat := c.Chat()
	if chat != nil {
		if chat.Title != "" {
			// group chat
			name = chat.Title
		} else {
			// private chat
			name = chat.Username
		}
	}

	err := tweetforward.Subscribe(models.Chat{ID: chat.ID, Name: name}, twuname)
	if err != nil {
		c.Send(fmt.Sprintf("[ERR] Fail to subscribe %s: %s", twuname, err.Error()))
		return errors.Wrap(err, twuname)
	}

	return c.Send("[OK] Subscribed to " + twuname)

}

func (t *Notifier) tweetList(c telebot.Context) error {
	chatid := c.Chat().ID

	subs, err := tweetforward.GetSubscriptions(chatid)
	if err != nil {
		return c.Send("[ERR] " + err.Error())
	}

	msg := new(strings.Builder)
	for _, sub := range subs {
		msg.WriteString(sub.Username + "\n")
	}

	return c.Send(msg.String())
}

func (t *Notifier) tweetUnSub(c telebot.Context) error {
	var twuname string

	args := c.Args()
	switch {
	case len(args) == 1:
		twuname = args[0]
	default:
		c.Send("usage: /tweetunsub <twitter username>")
		return ErrNotEnoughArguments
	}

	err := tweetforward.Unsubscribe(c.Chat().ID, twuname)
	if err != nil {
		c.Send(fmt.Sprintf("[ERR] Fail to unsubscribe %s: %s", twuname, err.Error()))
		return errors.Wrap(err, twuname)
	}

	return c.Send("[OK] Unsub to " + twuname)

}
