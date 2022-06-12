package main

import (
	"os"

	"github.com/0w0mewo/mcrc_tgbot/bot"
	_ "github.com/0w0mewo/mcrc_tgbot/bot/modules"
	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	defer persistent.DefaultDBConn.Close() // kill db connection when exit

	token := os.Getenv("TOKEN")
	if token == "" {
		logrus.Error("empty token")
		return
	}

	yoo, err := bot.NewBot(token)
	if err != nil {
		logrus.Errorln(err)
		return
	}

	yoo.Start()

	<-utils.WaitForSignal()

	yoo.Stop()

}
