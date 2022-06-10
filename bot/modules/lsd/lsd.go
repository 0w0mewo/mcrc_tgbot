package lsd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/service/linesticker"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

var ErrNotEnoughArguments = errors.New("not enough args")

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig("mod.lsd", cfg)

	pool, err := ants.NewPool(4)
	if err != nil {
		panic(err)
	}

	// load module
	m := &lineStickerDown{
		logger:  utils.NewLogger(),
		conf:    cfg,
		fetcher: linesticker.NewFetcher(context.Background(), http.DefaultClient),
		pool:    pool,
	}
	bot.ModRegister.RegistryMod(m)

}

type lineStickerDown struct {
	tgbot   *telebot.Bot
	conf    *lineStickerDownconf
	logger  *logrus.Logger
	fetcher *linesticker.Fetcher
	pool    *ants.Pool // for proccessing requested stickers packages
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
	return "mod.lsd"
}

func (ma *lineStickerDown) Stop(b *bot.Bot) {
	ma.running = false
	ma.pool.Release()

	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *lineStickerDown) Reload() {

}

func (ma *lineStickerDown) lsd(c telebot.Context) error {
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

	// pass request to workers
	err = ma.pool.Submit(func() {
		ma.downloadAndSend(c.Recipient(), packid, qqtrans)
	})
	if err != nil {
		ma.logger.Error("pool", err)
		return c.Send("fail to download package")
	}

	return c.Send(fmt.Sprintf("downloading sticker pack: %d", packid))

}

func (ma *lineStickerDown) downloadAndSend(respTo telebot.Recipient, packid int, qqTrans bool) {
	// download stickers
	data, err := ma.fetcher.SaveStickers(context.Background(), packid, qqTrans)
	if err != nil {
		ma.logger.Error("worker", err)
		ma.tgbot.Send(respTo, "unable to fetch stickers")
		return
	}

	// send
	zippedPack := &telebot.Document{
		File:     telebot.FromReader(bytes.NewReader(data)),
		FileName: fmt.Sprintf("%d.zip", packid),
	}

	ma.tgbot.Send(respTo, zippedPack, &telebot.SendOptions{
		DisableNotification: true,
	})

	return

}
