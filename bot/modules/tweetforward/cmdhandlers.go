package notifier

import (
	"fmt"
	"strings"

	"github.com/0w0mewo/mcrc_tgbot/bot/modules/commands/cmds"
	"github.com/0w0mewo/mcrc_tgbot/model"
	tweetforward "github.com/0w0mewo/mcrc_tgbot/service/tweetforward"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

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
		return cmds.ErrNotEnoughArguments
	}

	err := tweetforward.Subscribe(model.Chat{Id: c.Chat().ID, Name: c.Chat().Title}, twuname)
	if err != nil {
		c.Send(fmt.Sprintf("[ERR] Fail to subscribe %s: %s", twuname, err.Error()))
		return errors.Wrap(err, twuname)
	}

	return c.Send("[OK] Subscribed to " + twuname)

}

func (t *Notifier) tweetList(c telebot.Context) error {
	chat := c.Chat().ID

	subs, err := tweetforward.GetSubscriptions(chat)
	if err != nil {
		return c.Send("[ERR] " + err.Error())
	}

	msg := new(strings.Builder)
	for _, sub := range subs {
		msg.WriteString(sub.UserName + "\n")
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
		return cmds.ErrNotEnoughArguments
	}

	err := tweetforward.Unsubscribe(c.Chat().ID, twuname)
	if err != nil {
		c.Send(fmt.Sprintf("[ERR] Fail to unsubscribe %s: %s", twuname, err.Error()))
		return errors.Wrap(err, twuname)
	}

	return c.Send("[OK] Unsub to " + twuname)

}
