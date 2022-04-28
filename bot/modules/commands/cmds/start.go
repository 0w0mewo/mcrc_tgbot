package cmds

import (
	"gopkg.in/telebot.v3"
)

func Start(c telebot.Context) error {

	return c.Send("budong")
}
