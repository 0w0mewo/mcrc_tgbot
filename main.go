package main

import (
	"os"

	"github.com/0w0mewo/mcrc_tgbot/bot"
	_ "github.com/0w0mewo/mcrc_tgbot/bot/modules/commands"
	_ "github.com/0w0mewo/mcrc_tgbot/bot/modules/repeater"
	_ "github.com/0w0mewo/mcrc_tgbot/bot/modules/tweetforward"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
)

var retries = 0

func main() {
	// restart when crash
	defer func() {
		if err := recover(); err != nil {
			retries++
			if retries > 5 {
				return
			}

			logrus.Errorf("%s, retring...", err)
			main()

		}
	}()

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
