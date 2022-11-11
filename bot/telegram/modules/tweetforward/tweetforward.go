package notifier

import (
	tgbot "github.com/0w0mewo/mcrc_tgbot/bot"
	tweetforward "github.com/0w0mewo/mcrc_tgbot/service/tweetforward"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

const modname = "tg.mod.tweetforwarder"

func init() {
	// load module
	r := &Notifier{
		logger: utils.GetLogger().WithField("module", modname),
	}
	tgbot.TgModRegister.RegistryMod(r)
}

// tweet notifier
type Notifier struct {
	tgbot   *telebot.Bot
	logger  *logrus.Entry
	running bool
}

func (t *Notifier) Start(b *tgbot.TelegramBot) {
	if !t.running {
		t.tgbot = b.Bot()
		t.running = true
	}

	t.Reload()

	tweetforward.DoWhenSet(func(chatid int64, tweet string) error {
		_, err := t.tgbot.Send(&telebot.Chat{ID: chatid}, tweet)

		return err
	})

	tgbot.TgModRegister.AddTgEventHandler("/tweetsub", t.tweetSub)
	tgbot.TgModRegister.AddTgEventHandler("/tweetunsub", t.tweetUnSub)

	t.logger.Printf("%s loaded", t.Name())
}

func (t *Notifier) Name() string {
	return modname
}

func (t *Notifier) Stop(b *tgbot.TelegramBot) {
	t.running = false
	tweetforward.Shutdown()

	t.logger.Printf("%s unloaded", t.Name())
}

func (t *Notifier) Reload() {

}
