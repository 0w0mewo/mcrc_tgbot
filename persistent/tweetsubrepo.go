package persistent

import (
	"context"
	"database/sql"

	models "github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var _ ChatTweetSubRepo = &chatTweetSubSqlStorage{}

type chatTweetSubSqlStorage struct {
	conn *sql.DB
}

func NewChatTweetSubSqlStorage(db *sql.DB) *chatTweetSubSqlStorage {
	return &chatTweetSubSqlStorage{
		conn: db,
	}
}

func (this *chatTweetSubSqlStorage) Create(ctx context.Context, fromchat models.Chat, sub models.TweetUser) (err error) {
	tx, err := this.conn.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	// check if it has subscribed already
	exist, err := this.Exist(ctx, fromchat.ID, sub.ID)
	if err != nil {
		tx.Rollback()
		return
	}

	// it is
	if exist {
		tx.Rollback()
		return ErrExist
	}

	// upsert chat info
	err = fromchat.Upsert(ctx, tx, true, []string{models.ChatColumns.ID},
		boil.Whitelist(models.ChatColumns.Name), boil.Infer())
	if err != nil {
		tx.Rollback()
		return
	}

	// upsert twitter subscription
	err = sub.Upsert(ctx, tx, true, []string{models.TweetUserColumns.ID},
		boil.Whitelist(models.TweetUserColumns.Username), boil.Infer())
	if err != nil {
		tx.Rollback()
		return
	}

	// insert subscription
	newsub := models.ChatTweetSubscription{
		TweeterID: null.StringFrom(sub.ID),
		ChatID:    null.Int64From(fromchat.ID),
		LastTweet: "no tweet captured",
	}
	err = newsub.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer())
	if err != nil {
		tx.Rollback()
		return
	}

	return tx.Commit()
}

func (this *chatTweetSubSqlStorage) Remove(ctx context.Context, chatid int64, twuid string) error {
	exist, err := this.Exist(ctx, chatid, twuid)
	if err != nil {
		return err
	}

	if !exist {
		return ErrNotExist
	}

	return models.ChatTweetSubscriptions(
		chatAndTweeterFilter(chatid, twuid)...,
	).DeleteAll(ctx, this.conn)
}

func (this *chatTweetSubSqlStorage) Exist(ctx context.Context, chatid int64, twuid string) (bool, error) {
	return models.ChatTweetSubscriptions(
		qm.Select(models.ChatTweetSubscriptionColumns.ID),
		models.ChatTweetSubscriptionWhere.ChatID.EQ(null.Int64From(chatid)),
		models.ChatTweetSubscriptionWhere.TweeterID.EQ(null.StringFrom(twuid)),
	).Exists(ctx, this.conn)

}

func (this *chatTweetSubSqlStorage) GetAllChatSub(ctx context.Context, chatid int64) (models.ChatTweetSubscriptionSlice, error) {
	return models.ChatTweetSubscriptions(
		models.ChatTweetSubscriptionWhere.ChatID.EQ(null.Int64From(chatid))).
		All(ctx, this.conn)
}

func (this *chatTweetSubSqlStorage) GetAllChatIds(ctx context.Context) ([]int64, error) {
	res, err := models.ChatTweetSubscriptions(
		qm.Select(models.ChatTweetSubscriptionColumns.ChatID),
		qm.GroupBy(models.ChatTweetSubscriptionColumns.ChatID)).
		All(ctx, this.conn)
	if err != nil {
		return nil, err
	}

	chatids := make([]int64, 0, len(res))
	for _, c := range res {
		chatids = append(chatids, c.ChatID.Int64)
	}

	return chatids, nil
}

func (this *chatTweetSubSqlStorage) GetAllSubscribeeByChatId(ctx context.Context, chatid int64) (models.TweetUserSlice, error) {
	test, err := models.ChatTweetSubscriptions(
		qm.Select(models.ChatTweetSubscriptionColumns.TweeterID),
		models.ChatTweetSubscriptionWhere.ChatID.EQ(null.Int64From(chatid)),
	).All(ctx, this.conn)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(test))
	for _, t := range test {
		ids = append(ids, t.TweeterID.String)
	}

	return models.TweetUsers(models.TweetUserWhere.ID.IN(ids)).All(ctx, this.conn)
}

func (this *chatTweetSubSqlStorage) GetLastTweet(ctx context.Context, chatid int64, twuid string) (string, error) {
	res, err := models.ChatTweetSubscriptions(
		qm.Select(models.ChatTweetSubscriptionColumns.LastTweet),
		models.ChatTweetSubscriptionWhere.ChatID.EQ(null.Int64From(chatid)),
		models.ChatTweetSubscriptionWhere.TweeterID.EQ(null.StringFrom(twuid)),
	).One(ctx, this.conn)
	if err != nil {
		return "", err
	}

	return res.LastTweet, nil
}

func (this *chatTweetSubSqlStorage) UpdateLastTweet(ctx context.Context, chatid int64, twuid string, newtweet string) error {
	return models.ChatTweetSubscriptions(
		chatAndTweeterFilter(chatid, twuid)...,
	).UpdateAll(ctx, this.conn, models.M{models.ChatTweetSubscriptionColumns.LastTweet: newtweet})
}

func (this *chatTweetSubSqlStorage) GetTweeterOfChatSub(ctx context.Context, id int) (tu models.TweetUser, err error) {
	chatsub, err := models.ChatTweetSubscriptions(
		qm.Select(models.ChatTweetSubscriptionColumns.TweeterID),
		models.ChatTweetSubscriptionWhere.ID.EQ(int64(id)),
	).One(ctx, this.conn)
	if err != nil {
		return
	}

	res, err := models.TweetUsers(models.TweetUserWhere.ID.EQ(chatsub.TweeterID.String)).One(ctx, this.conn)
	if err != nil {
		return
	}

	tu.ID = res.ID
	tu.Username = res.Username

	return
}

func (this *chatTweetSubSqlStorage) Close() error {
	return nil
}

func chatAndTweeterFilter(chatid int64, twuid string) []qm.QueryMod {
	return []qm.QueryMod{
		models.ChatTweetSubscriptionWhere.ChatID.EQ(null.Int64From(chatid)),
		models.ChatTweetSubscriptionWhere.TweeterID.EQ(null.StringFrom(twuid)),
	}
}
