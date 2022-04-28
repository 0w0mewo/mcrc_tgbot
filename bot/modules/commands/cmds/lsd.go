package cmds

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/service/linesticker"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"gopkg.in/telebot.v3"
)

func Lsd(c telebot.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var packid int
	var qqtrans bool
	var err error

	// expect /lsd <packid> [<qqtrans>]
	args := c.Args()
	switch {
	case len(args) == 1:
		packid, err = strconv.Atoi(args[0])
		if err != nil {
			c.Send("invalid pack id")

			return err
		}
	case len(args) == 2:
		packid, err = strconv.Atoi(args[0])
		if err != nil {
			c.Send("invalid pack id")

			return err
		}
		qqtrans = utils.StringToBoolean(args[1])
	default:
		c.Send("usage: /lsd <LINE sticker package id>")
		return ErrNotEnoughArguments
	}

	// download stickers
	fetcher := linesticker.NewFetcher(ctx, http.DefaultClient)
	data, err := fetcher.SaveStickers(ctx, packid, qqtrans)
	if err != nil {
		return err
	}

	// send
	zippedPack := &telebot.Document{
		File:     telebot.FromReader(bytes.NewReader(data)),
		FileName: fmt.Sprintf("%d.zip", packid),
	}

	return c.Send(zippedPack, &telebot.SendOptions{
		DisableNotification: true,
	})
}
