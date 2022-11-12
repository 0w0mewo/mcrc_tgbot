package avatar

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"time"

	dcbot "github.com/0w0mewo/mcrc_tgbot/bot"
	"github.com/0w0mewo/mcrc_tgbot/config"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const modname = "dc.mod.avatar"

func init() {
	// load config and regsiter to manager
	cfg := ConfigFrom(config.GetConfigFile())
	config.RegisterModuleConfig(modname, cfg)

	// load module
	m := &avatar{
		logger:    utils.GetLogger().WithField("module", modname),
		conf:      cfg,
		scheduled: utils.GetDefaultScheduledTasksGrp(),
	}
	dcbot.DcModRegister.RegistryMod(m)

}

type avatar struct {
	dcbot     *discordgo.Session
	conf      *avatarConf
	logger    *logrus.Entry
	scheduled *utils.ScheduledTaskGroup
	running   bool
}

func (ma *avatar) Start(b *dcbot.DiscordBot) {
	if !ma.running {
		ma.dcbot = b.Bot()
		ma.running = true
	}

	ma.Reload()

	ma.logger.Infof("watch to %s", ma.conf.syncToUserId)
	ma.scheduled.AddPerodical(5*time.Minute, func() error {
		return ma.avatard()
	})

	ma.logger.Printf("%s loaded", ma.Name())
}

func (ma *avatar) Name() string {
	return modname
}

func (ma *avatar) Stop(b *dcbot.DiscordBot) {
	ma.running = false
	ma.scheduled.WaitAndStop()

	ma.logger.Printf("%s unloaded", ma.Name())
}

func (ma *avatar) Reload() {

}

func (ma *avatar) avatard() (err error) {
	watchTo := ma.conf.syncToUserId

	expected, username, err := getUserAvatar(ma.dcbot, watchTo)
	if err != nil {
		ma.logger.Error(err)
		return
	}

	current, _, err := getUserAvatar(ma.dcbot, ma.dcbot.State.User.ID)
	if err != nil {
		ma.logger.Error(err)
		return
	}

	dist, pdiff, err := utils.CompareTwoImage(expected, current)
	if err != nil {
		ma.logger.Error(err)
		return
	}

	ma.logger.Debugf("distance to expected avatar: %f, pixel difference: %f %%", dist, pdiff*100.0)

	if pdiff > 0.5 {
		ma.logger.Println("avatar change detected")

		imgbuf := bytes.NewBuffer(nil)
		err = png.Encode(imgbuf, expected)
		if err != nil {
			ma.logger.Error(err)
			return
		}

		base64img := base64.StdEncoding.EncodeToString(imgbuf.Bytes())
		avatar := fmt.Sprintf("data:image/png;base64,%s", base64img)

		_, err = ma.dcbot.UserUpdate(username, avatar)
		if err != nil {
			ma.logger.Error(err)
			return
		}

		return
	}

	return

}

func getUserAvatar(s *discordgo.Session, userId string) (img image.Image, username string, err error) {
	u, err := s.User(userId)
	if err != nil {
		return
	}

	img, err = s.UserAvatarDecode(u)
	username = u.Username
	return

}
