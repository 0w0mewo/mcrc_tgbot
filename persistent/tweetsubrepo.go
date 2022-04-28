package persistent

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent/chat"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent/chattweetsubscription"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent/predicate"
	"github.com/0w0mewo/mcrc_tgbot/persistent/ent/tweetuser"
)

var _ ChatTweetSubRepo = &chatTweetSubSqlStorage{}

type chatTweetSubSqlStorage struct {
	dbconn *ent.Client
}

func NewChatTweetSubSqlStorage(dbconn *ent.Client) ChatTweetSubRepo {
	ret := &chatTweetSubSqlStorage{
		dbconn: dbconn,
	}

	return ret
}

func (cts *chatTweetSubSqlStorage) Create(ctx context.Context, fromchat model.Chat, sub model.TweetUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// check existence first
	exist, err := cts.Exist(ctx, fromchat.Id, sub.Id)
	if err != nil {
		return err
	}
	if exist {
		return ErrExist
	}

	// start a transaction since there are two table related together
	tx, err := cts.dbconn.Tx(ctx)
	if err != nil {
		return err
	}

	// add tweeter user detail
	err = tx.TweetUser.Create().
		SetID(sub.Id).
		SetUsername(sub.UserName).
		OnConflict(
			sql.ConflictColumns(tweetuser.FieldID),
		).
		UpdateNewValues().
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err

	}

	// work around of FOREIGN KEY constraint fail
	err = tx.Chat.Create().
		SetID(fromchat.Id).
		SetName(fromchat.Name).
		OnConflict(
			sql.ConflictColumns(chat.FieldID),
		).UpdateNewValues().Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// create chat-tweeter subscription record, rollback if failed to create
	err = tx.ChatTweetSubscription.Create().
		SetChatID(fromchat.Id).
		SetLastTweet(sub.LastTweet).
		SetTweeterID(sub.Id).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()

}

func (cts *chatTweetSubSqlStorage) Remove(ctx context.Context, chatid int64, twuid string) error {
	exist, err := cts.Exist(ctx, chatid, twuid)
	if err != nil {
		return err
	}
	if !exist {
		return ErrNotExist
	}

	_, err = cts.dbconn.ChatTweetSubscription.Delete().
		Where(cts.matchTwetterIdAndChatId(chatid, twuid)...).
		Exec(ctx)

	return err
}

func (cts *chatTweetSubSqlStorage) Exist(ctx context.Context, chatid int64, twuid string) (bool, error) {
	return cts.dbconn.ChatTweetSubscription.Query().
		Where(cts.matchTwetterIdAndChatId(chatid, twuid)...).
		Exist(ctx)
}

func (cts *chatTweetSubSqlStorage) GetAllChatSub(ctx context.Context,
	chatid int64) ([]*model.ChatTweetSubscription, error) {
	subs, err := cts.dbconn.ChatTweetSubscription.Query().
		Where(chattweetsubscription.ChatID(chatid)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]*model.ChatTweetSubscription, 0, len(subs))

	for _, cs := range subs {
		res = append(res, &model.ChatTweetSubscription{
			Id:        cs.ID,
			ChatId:    cs.ChatID,
			TweeterId: cs.TweeterID,
			LastTweet: cs.LastTweet,
		})
	}

	return res, nil
}

func (cts *chatTweetSubSqlStorage) GetAllChat(ctx context.Context) ([]int64, error) {
	res, err := cts.dbconn.ChatTweetSubscription.Query().
		GroupBy(chattweetsubscription.FieldChatID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}

	dbChatid := make([]int64, 0, len(res))
	for _, c := range res {
		dbChatid = append(dbChatid, int64(c))
	}

	return dbChatid, nil
}

func (cts *chatTweetSubSqlStorage) GetLastTweet(ctx context.Context, chatid int64, twuid string) (string, error) {
	res, err := cts.dbconn.ChatTweetSubscription.Query().
		Where(cts.matchTwetterIdAndChatId(chatid, twuid)...).
		First(ctx)
	if err != nil {
		return "", err
	}

	return res.LastTweet, nil
}

func (cts *chatTweetSubSqlStorage) UpdateLastTweet(ctx context.Context, chatid int64, twuid string, newtweet string) error {
	_, err := cts.dbconn.ChatTweetSubscription.Update().
		Where(cts.matchTwetterIdAndChatId(chatid, twuid)...).
		SetLastTweet(newtweet).
		Save(ctx)

	return err
}

func (cts *chatTweetSubSqlStorage) GetTweeterOfChatSub(ctx context.Context, id int) (*model.TweetUser, error) {
	chatsub, err := cts.dbconn.ChatTweetSubscription.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	twu, err := cts.dbconn.ChatTweetSubscription.QuerySubscribedTweeter(chatsub).
		First(ctx)
	if err != nil {
		return nil, err
	}

	return &model.TweetUser{
		Id:       twu.ID,
		UserName: twu.Username,
	}, nil

}

func (cts *chatTweetSubSqlStorage) GetAllSubscribeeByChatId(ctx context.Context, chatid int64) ([]*model.TweetUser, error) {
	res, err := cts.dbconn.ChatTweetSubscription.Query().Where(chattweetsubscription.ChatID(chatid)).QuerySubscribedTweeter().All(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]*model.TweetUser, 0, len(res))

	for _, twu := range res {
		ret = append(ret, &model.TweetUser{Id: twu.ID, UserName: twu.Username})
	}

	return ret, nil
}

func (cts *chatTweetSubSqlStorage) Close() error {
	return cts.dbconn.Close()
}

func (cts *chatTweetSubSqlStorage) matchTwetterIdAndChatId(chatid int64, twuid string) []predicate.ChatTweetSubscription {
	return []predicate.ChatTweetSubscription{
		chattweetsubscription.ChatID(chatid),
		chattweetsubscription.TweeterID(twuid),
	}
}
