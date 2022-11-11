package main

import (
	"os"

	"github.com/0w0mewo/mcrc_tgbot/bot"
	_ "github.com/0w0mewo/mcrc_tgbot/bot/discord/modules"
	_ "github.com/0w0mewo/mcrc_tgbot/bot/telegram/modules"
	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	defer persistent.DefaultDBConn.Close() // kill db connection when exit

	tgtoken := os.Getenv("TGTOKEN")
	if tgtoken == "" {
		logrus.Error("empty token")
		return
	}

	dctoken := os.Getenv("DCTOKEN")
	if dctoken == "" {
		logrus.Error("empty token")
		return
	}

	yoo, err := bot.NewTelegramBot(tgtoken)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	yoo2, err := bot.NewDiscordBot(dctoken)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	yoo.Start()
	yoo2.Start()

	<-utils.WaitForSignal()

	yoo.Stop()
	yoo2.Stop()
}
