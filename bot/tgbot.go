package bot

import (
	"time"

	"github.com/0w0mewo/mcrc_tgbot/bot/telegram"
	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

type TelegramBot struct {
	tgbot             *telebot.Bot
	logger            *logrus.Entry
	msgrepo           persistent.StoredTeleMsgRepo
	eventHandlersPool *ants.Pool
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	logger := utils.GetLogger().WithField("service", "tgbot")

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

	return &TelegramBot{
		tgbot:             tgbot,
		logger:            logger,
		msgrepo:           repo,
		eventHandlersPool: pool,
	}, nil
}

func (b *TelegramBot) Start() {
	// load middlewares
	b.tgbot.Use(telegram.BypassBotMessage, telegram.BypassSelfMessage, telegram.StoreGrpMessage(b.msgrepo))

	// load modules
	mods := TgModRegister.GetModules()
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

func (b *TelegramBot) Stop() {
	defer func() {
		b.tgbot.Stop()
		b.logger.Infof("stopped")
	}()
	// kill all modules
	mods := TgModRegister.GetModules()
	for _, m := range mods {
		m.Stop(b)
	}

	// kill repo
	b.msgrepo.Close()

	// kill event handlers pool
	b.eventHandlersPool.Release()

}

func (b *TelegramBot) Bot() *telebot.Bot {
	return b.tgbot
}

func (b *TelegramBot) processHandlers(ev string) telebot.HandlerFunc {
	return func(c telebot.Context) (err error) {
		// process all registered handlers
		handlers := TgModRegister.GetTgEventHandlers(ev)

		for _, handler := range handlers {
			b.eventHandlersPool.Submit(func() {
				handler(c)
			})
		}

		return
	}
}
