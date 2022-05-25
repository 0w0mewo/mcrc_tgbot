package notifier

import (
	"github.com/0w0mewo/mcrc_tgbot/bot"
	tweetforward "github.com/0w0mewo/mcrc_tgbot/service/tweetforward"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

func init() {
	// load module
	r := &Notifier{
		logger: utils.NewLogger(),
	}
	bot.ModRegister.RegistryMod(r)
}

// group message random repeater
type Notifier struct {
	tgbot *telebot.Bot

	logger  *logrus.Logger
	running bool
}

func (t *Notifier) Start(b *bot.Bot) {
	if !t.running {
		t.tgbot = b.Bot()
		t.running = true
	}

	t.Reload()

	tweetforward.DoWhenSet(func(chatid int64, tweet string) error {
		_, err := t.tgbot.Send(&telebot.Chat{ID: chatid}, tweet)

		return err
	})

	t.logger.Printf("%s loaded", t.Name())
}

func (t *Notifier) Name() string {
	return "mod.tweetforwarder"
}

func (t *Notifier) Stop(b *bot.Bot) {
	t.running = false
	tweetforward.Shutdown()

	t.logger.Printf("%s unloaded", t.Name())
}

func (t *Notifier) Reload() {
	t.tgbot.Handle("/tweetsub", t.tweetSub)
	t.tgbot.Handle("/tweetunsub", t.tweetUnSub)

}
