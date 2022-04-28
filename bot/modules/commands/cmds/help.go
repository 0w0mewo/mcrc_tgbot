package cmds

import "gopkg.in/telebot.v3"

func Help(c telebot.Context) error {
	return c.Send(`
		/lsd <pack id> : fetch LINE stickers package
		/tweetsub <twitter username> : subscribe tweeter
		/tweetunsub <twitter username> : unsubscribe tweeter 

	`)
}
