package persistent

import (
	"context"
	"database/sql"
	"time"

	models "github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	conn *sql.DB
}

func NewTeleMsgSqlStorage(db *sql.DB) *teleMsgRecSqlStorage {
	return &teleMsgRecSqlStorage{
		conn: db,
	}

}

func (this *teleMsgRecSqlStorage) StoreMsg(ctx context.Context, chatid int64, chatname string, senderid int64, sendername string, msg []byte, msgType int, timestamp time.Time) (err error) {
	tx, err := this.conn.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	// upsert chat info
	chat := models.Chat{
		ID:   chatid,
		Name: chatname,
	}
	err = chat.Upsert(ctx, tx, true,
		[]string{models.ChatColumns.ID},
		boil.Whitelist(models.ChatColumns.Name),
		boil.Infer())
	if err != nil {
		tx.Rollback()
		return
	}

	// upsert sender info
	sender := models.Sender{
		ID:       senderid,
		Username: sendername,
	}
	err = sender.Upsert(ctx, tx, true,
		[]string{models.SenderColumns.ID},
		boil.Whitelist(models.SenderColumns.Username),
		boil.Infer())
	if err != nil {
		tx.Rollback()
		return
	}

	// insert message
	newrow := models.Message{
		MSG:       msg,
		Timestamp: timestamp,
		ChatID:    null.Int64From(chatid),
		SenderID:  null.Int64From(senderid),
		Type:      int64(msgType),
	}
	err = newrow.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()) // no way it would conflicts
	if err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()

}

func (this *teleMsgRecSqlStorage) GetChatMessages(ctx context.Context, chatid int64, page int, pagesize int) (models.MessageSlice, error) {
	return models.Messages(
		models.MessageWhere.ChatID.EQ(null.Int64From(chatid)),
		qm.Limit(pagesize),
		qm.Offset(page),
	).All(ctx, this.conn)
}

func (this *teleMsgRecSqlStorage) Count(ctx context.Context) (int, error) {
	cnt, err := models.Messages().Count(ctx, this.conn)

	return int(cnt), err
}

func (this *teleMsgRecSqlStorage) Close() error {
	return nil
}
