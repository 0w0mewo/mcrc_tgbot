package bot

import (
	"time"

	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	tgbot             *telebot.Bot
	logger            *logrus.Entry
	msgrepo           persistent.StoredTeleMsgRepo
	eventHandlersPool *ants.Pool
}

func NewBot(token string) (*Bot, error) {
	logger := utils.NewLogger().WithField("service", "tgbot")

	// bot init
	tgbot, err := telebot.NewBot(telebot.Settings{
		Token:   token,
		Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
		Verbose: false,
		OnError: func(err error, ctx telebot.Context) {
			chat := ctx.Chat().Title
			sender := ctx.Sender().Username
			updateid := ctx.Update().ID

			logger.WithField("chat", chat).
				WithField("sender", sender).
				WithField("update id", updateid).
				Error(err)
		},
	})
	if err != nil {
		return nil, err
	}

	// messages storage, DB session
	repo := persistent.NewTeleMsgSqlStorage(persistent.DefaultDBConn)

	// pool of workers for event handlers
	pool, err := ants.NewPool(10)
	if err != nil {
		return nil, err
	}

	return &Bot{
		tgbot:             tgbot,
		logger:            logger,
		msgrepo:           repo,
		eventHandlersPool: pool,
	}, nil
}

func (b *Bot) Start() {
	// load middlewares
	b.tgbot.Use(BypassBotMessage, BypassSelfMessage, StoreGrpMessage(b.msgrepo))

	// load modules
	mods := ModRegister.GetModules()
	for _, m := range mods {
		m.Start(b)
	}

	// registry handlers
	for _, ev := range listenTo {
		b.tgbot.Handle(ev, b.processHandlers(ev))
	}

	// start
	go b.tgbot.Start()

}

func (b *Bot) Stop() {
	// kill all modules
	mods := ModRegister.GetModules()
	for _, m := range mods {
		m.Stop(b)
	}

	// kill repo
	b.msgrepo.Close()

	// kill event handlers pool
	b.eventHandlersPool.Release()

	// stop tgbot
	b.tgbot.Stop()
}

func (b *Bot) Bot() *telebot.Bot {
	return b.tgbot
}

func (b *Bot) processHandlers(ev string) telebot.HandlerFunc {
	return func(c telebot.Context) (err error) {
		// process all registered handlers
		handlers := ModRegister.GetTgEventHandlers(ev)

		for _, handler := range handlers {
			// go func(handler telebot.HandlerFunc) {
			// 	err = handler(c)
			// }(handler)

			b.eventHandlersPool.Submit(func() {
				handler(c)
			})
		}

		return
	}
}
