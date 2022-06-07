package lsd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/service/linesticker"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

var ErrNotEnoughArguments = errors.New("not enough args")

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.Config.GetConfigFile())
	config.Config.RegisterModuleConfig("mod.lineStickerDown", cfg)

	// load module
	m := &lineStickerDown{
		logger: utils.NewLogger(),
		conf:   cfg,
	}
	bot.ModRegister.RegistryMod(m)

}

type lineStickerDown struct {
	tgbot   *telebot.Bot
	conf    *lineStickerDownconf
	logger  *logrus.Logger
	running bool
}

func (ma *lineStickerDown) Start(b *bot.Bot) {
	if !ma.running {
		ma.tgbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	bot.ModRegister.AddTgEventHandler("/lsd", ma.lsd)

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *lineStickerDown) Name() string {
	return "mod.lineStickerDown"
}

func (ma *lineStickerDown) Stop(b *bot.Bot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *lineStickerDown) Reload() {

}

func (ma *lineStickerDown) lsd(c telebot.Context) error {
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
