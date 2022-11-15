package lsd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	dcbot "github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/service/linesticker"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

var ErrNotEnoughArguments = errors.New("not enough args")

const modname = "dc.mod.lsd"
const cmd = "lsd"

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig(modname, cfg)

	pool, err := ants.NewPool(4)
	if err != nil {
		panic(err)
	}

	// load module
	m := &lsd{
		logger:  utils.GetLogger().WithField("module", modname),
		conf:    cfg,
		fetcher: linesticker.NewFetcher(context.Background(), http.DefaultClient),
		pool:    pool,
	}
	dcbot.DcModRegister.RegistryMod(m)

}

type lsd struct {
	dcbot   *discordgo.Session
	conf    *lsdConf
	fetcher *linesticker.Fetcher
	pool    *ants.Pool // for proccessing requested stickers packages
	logger  *logrus.Entry
	running bool
}

func (ma *lsd) Start(b *dcbot.DiscordBot) {
	if !ma.running {
		ma.dcbot = b.Bot()
		ma.running = true
	}

	// non-slash command
	dcbot.DcModRegister.AddDcMesgHandler(
		dcbot.WrappedDiscordCmdHandler(cmd, ma.lsd))

	// slash command
	desc := dcbot.NewDiscordSlashCmdEntry(&discordgo.ApplicationCommand{
		Name:        cmd,
		Description: "line stickers downloader",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "packid",
				Description: "stickers package id",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "qqtrans",
				Description: "whether the stickers convert to gif format",
				Required:    false,
			},
		},
	}, ma.slashlsd)
	dcbot.DcModRegister.AddDcSlashCmdHandler(cmd, desc)

	// reload
	ma.Reload()

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *lsd) Name() string {
	return modname
}

func (ma *lsd) Stop(b *dcbot.DiscordBot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *lsd) Reload() {

}

func (ma *lsd) slashlsd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_params := i.ApplicationCommandData().Options
	parmas := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(_params))
	for _, opt := range _params {
		parmas[opt.Name] = opt
	}

	packid := parmas["packid"].IntValue()
	_qqtrans := parmas["qqtrans"]
	var qqtrans bool
	if _qqtrans == nil {
		qqtrans = true
	} else {
		qqtrans = _qqtrans.BoolValue()
	}

	// pass request to workers
	err := ma.pool.Submit(func() {
		ma.downloadAndSend(i.ChannelID, int(packid), qqtrans)
	})
	if err != nil {
		ma.logger.Error("pool", err)
		utils.DiscordSlashCmdRespString(s, i.Interaction, "fail to download package")
		return
	}

	utils.DiscordSlashCmdRespString(s, i.Interaction, fmt.Sprintf("downloading sticker pack: %d", packid))
}

func (ma *lsd) lsd(s *discordgo.Session, m *discordgo.MessageCreate) {
	var packid int
	var qqtrans bool
	var err error

	senderChannel := m.ChannelID

	// expect >lsd <packid> [<qqtrans>]
	args := strings.Split(m.Content, " ")
	args = args[1:]

	switch {
	case len(args) == 1:
		packid, err = strconv.Atoi(args[0])
		if err != nil {
			s.ChannelMessageSend(senderChannel, "invalid pack id")
			return
		}
	case len(args) == 2:
		packid, err = strconv.Atoi(args[0])
		if err != nil {
			s.ChannelMessageSend(senderChannel, "invalid pack id")
			return
		}
		qqtrans = utils.StringToBoolean(args[1])
	default:
		s.ChannelMessageSend(senderChannel, "usage: >lsd <LINE sticker package id>")
		return
	}

	// pass request to workers
	err = ma.pool.Submit(func() {
		ma.downloadAndSend(senderChannel, packid, qqtrans)
	})
	if err != nil {
		ma.logger.Error("pool", err)
		s.ChannelMessageSend(senderChannel, "fail to download package")
		return
	}

	s.ChannelMessageSend(senderChannel, fmt.Sprintf("downloading sticker pack: %d", packid))

}

func (ma *lsd) downloadAndSend(respTo string, packid int, qqTrans bool) {
	// download stickers
	data, err := ma.fetcher.SaveStickers(context.Background(), packid, qqTrans)
	if err != nil {
		ma.logger.Error("worker: ", err)
		ma.dcbot.ChannelMessageSend(respTo, "fail to download")
		return
	}

	zippedPack := &discordgo.MessageSend{
		File: &discordgo.File{
			Name:   fmt.Sprintf("%d.zip", packid),
			Reader: bytes.NewReader(data),
		},
	}

	_, err = ma.dcbot.ChannelMessageSendComplex(respTo, zippedPack)
	if err != nil {
		ma.logger.Error(err)
		ma.dcbot.ChannelMessageSend(respTo, err.Error())
	}

}
