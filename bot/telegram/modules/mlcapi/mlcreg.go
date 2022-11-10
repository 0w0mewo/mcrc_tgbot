package mlcapi

import (
	"fmt"
	"time"

	tgbot "github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/service/mlcapi"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

const modname = "tg.mod.mlcapi"

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig(modname, cfg)

	// load module
	m := &MlcApi{
		logger: utils.NewLogger(),
		conf:   cfg,
		mlc:    mlcapi.NewMlcApiClient("", ""),
	}
	tgbot.TgModRegister.RegistryMod(m)

}

type MlcApi struct {
	tgbot   *telebot.Bot
	conf    *mlcApiConf
	logger  *logrus.Logger
	mlc     *mlcapi.MlcApiClient
	running bool
}

func (ma *MlcApi) Start(b *tgbot.TelegramBot) {
	if !ma.running {
		ma.tgbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	tgbot.TgModRegister.AddTgEventHandler("/mlcreg", ma.mlcreg)

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *MlcApi) Name() string {
	return modname
}

func (ma *MlcApi) Stop(b *tgbot.TelegramBot) {
	ma.running = false
	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *MlcApi) Reload() {
	ma.mlc.SetManagerToken(ma.conf.managerToken)
	ma.mlc.SetStaffToken(ma.conf.staffToken)

}

func (ma *MlcApi) mlcreg(c telebot.Context) error {
	if sender := c.Sender(); sender != nil {
		if !utils.IsInSlice(ma.conf.lockto, sender.Username) {
			return c.Send("you are not my master!!!")
		}
	}

	args := c.Args()

	if len(args) < 2 {
		return c.Send("usage: /mlcreg [register | deleteuser | changepw | setuserrak | unban | ban | resethwid | getinfo] <username> [<password> | <ban reason> | <rank>] [<rank>]")

	}

	res, err := mlcapi.ApiCall(ma.mlc, args[0], args[1:]...)
	if err != nil {
		return c.Send("ERROR: " + err.Error())
	}

	if res != nil {
		resp := res.(*mlcapi.UserInfo)
		return c.Send(fmt.Sprintf("username: %s, state: %s, last logined IP: %s, last logined at: %s, rank: %s",
			resp.UserName, resp.State, resp.LastIP, time.UnixMilli(resp.LastLoginTime), resp.Rank))
	}

	return c.Send("OK")

}
