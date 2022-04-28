package bot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"gopkg.in/telebot.v3"
)

func StoreGrpMessage(repo persistent.StoredTeleMsgRepo) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			msg := c.Message()
			// sender := msgrec.Sender{ID: msg.Sender.ID, UserName: msg.Sender.Username}
			// chat := msgrec.Chat{ID: msg.Chat.ID, Name: msg.Chat.Title}

			if !msg.FromGroup() {
				goto donext
			}

			switch {
			// text message
			case msg.Text != "":
				if strings.HasPrefix(msg.Text, "/") {
					goto donext
				}

				err := repo.StoreMsg(context.Background(),
					msg.Chat.ID, msg.Chat.Title, // chat
					msg.Sender.ID, msg.Sender.Username, // sender
					[]byte(msg.Text), persistent.MEDIA_TEXT, msg.Time()) // message body
				if err != nil {
					return err
				}

				// media
			case msg.Media() != nil:
				media := msg.Media()
				mediaData := &bytes.Buffer{}
				mediaType := persistent.MediaTypeFromString(media.MediaType())

				// fetch media file from TG server
				file, err := c.Bot().File(media.MediaFile())
				if err != nil {
					return err
				}
				defer file.Close()
				io.Copy(mediaData, file)

				// FIXME: media data may be empty for unknown reason
				if len(mediaData.Bytes()) == 0 {
					return fmt.Errorf("media file is empty: %s", media.MediaType())
				}

				// store it to database
				err = repo.StoreMsg(context.Background(),
					msg.Chat.ID, msg.Chat.Title, // chat
					msg.Sender.ID, msg.Sender.Username, // sender
					mediaData.Bytes(), mediaType, msg.Time()) //message body
				if err != nil {
					return err
				}

			}
		donext:
			return next(c)
		}
	}
}

func MessageCounter(cnter *utils.Counter, randlimit *utils.RandomMap) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			msg := c.Message()
			chatId := strconv.FormatInt(msg.Chat.ID, 10)

			// load random message count limit of a chat
			randlimit.Get(chatId)

			cnter.Inc(chatId)

			// if the current amount of message count over the limit
			if cnter.Get(chatId) > randlimit.Get(chatId) {
				randlimit.Generate(chatId)
				cnter.Reset(chatId)

				return next(c)
			}

			return nil
		}
	}
}

func BypassBotMessage(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		sender := c.Sender()
		if sender.IsBot {
			return nil
		}

		return next(c)

	}
}

func BypassSelfMessage(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		sender := c.Sender()
		if sender.ID == c.Bot().Me.ID {
			return nil
		}

		return next(c)
	}
}
