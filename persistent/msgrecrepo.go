package persistent

import (
	"context"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent/message"
)

const (
	MEDIA_VOICE = iota
	MEDIA_VIDEO
	MEDIA_AUDIO
	MEDIA_TEXT
	MEDIA_FILE
	MEDIA_PHOTO
	MEDIA_UNKONWN
)

func MediaTypeFromString(s string) int {
	switch s {
	case "document":
		return MEDIA_FILE
	case "photo":
		return MEDIA_PHOTO
	case "audio":
		return MEDIA_AUDIO
	case "video":
		return MEDIA_VIDEO
	case "voice":
		return MEDIA_VOICE
	default:
		return MEDIA_UNKONWN
	}
}

var _ StoredTeleMsgRepo = &teleMsgRecSqlStorage{}

type teleMsgRecSqlStorage struct {
	dbconn *ent.Client
}

func NewTeleMsgSqlStorage(dbclient *ent.Client) StoredTeleMsgRepo {
	ret := &teleMsgRecSqlStorage{
		dbconn: dbclient,
	}

	return ret
}

func (mrs *teleMsgRecSqlStorage) Close() error {
	return mrs.dbconn.Close()
}

func (mrs *teleMsgRecSqlStorage) GetChatMessages(ctx context.Context, chatid int64,
	page int, pagesize int) ([]*model.TeleMsg, error) {
	// get limited number of messages
	pagedMsgs, err := mrs.dbconn.Message.Query().
		Where(message.ChatID(chatid)).
		Limit(pagesize).Offset(page).
		All(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*model.TeleMsg, 0, len(pagedMsgs))

	for _, msg := range pagedMsgs {
		chat, err := msg.QueryFromChat().First(ctx)
		if err != nil {
			continue
		}
		sender, err := msg.QueryFromSender().First(ctx)
		if err != nil {
			continue
		}

		res = append(res, &model.TeleMsg{
			ChatName:       chat.Name,
			SenderUserName: sender.Username,
			Type:           msg.Type,
			Message:        msg.Msg,
			TimeStamp:      msg.Timestamp,
		})
	}

	return res, nil
}

func (mrs *teleMsgRecSqlStorage) StoreMsg(ctx context.Context, chatid int64, chatname string,
	senderid int64, sendername string, msg []byte, msgType int, timestamp time.Time) error {
	tx, err := mrs.dbconn.Tx(ctx)
	if err != nil {
		return err
	}

	// sender
	_, err = tx.Sender.Create().
		SetID(senderid).
		SetUsername(sendername).
		Save(ctx)
	if err != nil {
		// sender info had already recorded
		if ent.IsConstraintError(err) {
			goto recordmsg
		} else {
			tx.Rollback()
			return err
		}
	}

	// chat
	_, err = tx.Chat.Create().
		SetID(chatid).
		SetName(chatname).
		Save(ctx)
	if err != nil {
		// chat info had already recorded
		if ent.IsConstraintError(err) {
			goto recordmsg
		} else {
			tx.Rollback()
			return err
		}
	}

recordmsg:
	_, err = tx.Message.Create().
		SetChatID(chatid).
		SetSenderID(senderid).
		SetMsg(msg).
		SetTimestamp(timestamp).
		SetType(msgType).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (mrs *teleMsgRecSqlStorage) Count(ctx context.Context) (int, error) {
	return mrs.dbconn.Message.Query().Count(ctx)
}
