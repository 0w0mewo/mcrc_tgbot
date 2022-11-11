package tweetforward

import (
	"context"
	"errors"
	"sync"

	constant "github.com/0w0mewo/mcrc_tgbot/const"
	models "github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/agoalofalife/event"
	"github.com/sirupsen/logrus"
)

type tweetForwarder struct {
	repo           persistent.ChatTweetSubRepo
	logger         *logrus.Entry
	evhub          *event.Dispatcher
	scheduledTasks *utils.ScheduledTaskGroup
	wg             *sync.WaitGroup
}

func newTweetForwarder() *tweetForwarder {
	tf := &tweetForwarder{
		repo:           persistent.NewChatTweetSubSqlStorage(persistent.DefaultDBConn),
		logger:         utils.GetLogger().WithField("service", "tweetlistener"),
		evhub:          event.New(),
		scheduledTasks: utils.NewScheduledTaskGroup("tweet_update"),
		wg:             &sync.WaitGroup{},
	}

	// register periodical poller
	tf.scheduledTasks.AddPerodical(constant.TWEET_REFRESH_INTERVAL, tf.poll)

	return tf
}

func (tl *tweetForwarder) GetSubscriptions(chatid int64) (models.TweetUserSlice, error) {
	return tl.repo.GetAllSubscribeeByChatId(context.Background(), chatid)
}

func (tl *tweetForwarder) Subscribe(chat models.Chat, tweeter string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DB_CRUD_TIMEOUT)
	defer cancel()

	twu, lasttweet, err := ScrapLastTweet(tweeter)
	if err != nil {
		return err
	}

	err = tl.repo.Create(ctx, chat, *twu)
	if errors.Is(err, persistent.ErrExist) {
		return ErrSubscribed
	}

	return tl.repo.UpdateLastTweet(ctx, chat.ID, twu.ID, lasttweet)

}

func (tl *tweetForwarder) Unsubscribe(chatid int64, tweeter string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DB_CRUD_TIMEOUT)
	defer cancel()

	twu, _, err := ScrapLastTweet(tweeter)
	if err != nil {
		return err
	}

	err = tl.repo.Remove(ctx, chatid, twu.ID)
	if errors.Is(err, persistent.ErrNotExist) {
		return ErrNotSubscribed
	}

	return err
}

func (tl *tweetForwarder) DoWhenSet(cb any) {
	tl.evhub.Add("newtweet", cb)
}

func (tl *tweetForwarder) Shutdown() {
	err := tl.repo.Close()
	if err != nil {
		tl.logger.Error(err)
	}

	tl.scheduledTasks.WaitAndStop()

	tl.wg.Wait()
}

func (tl *tweetForwarder) poll() error {
	// get all chats which subscribed various tweeters
	chats, err := tl.repo.GetAllChatIds(context.Background())
	if err != nil {
		tl.logger.Errorln(err)
		return err
	}

	// go through each chat
	for _, chatid := range chats {
		tl.wg.Add(1)
		go tl.updateChatSubs(chatid)
	}

	return nil
}

func (tl *tweetForwarder) updateChatSubs(chatid int64) error {
	defer tl.wg.Done()

	subs, err := tl.repo.GetAllChatSub(context.Background(), chatid)
	if err != nil {
		tl.logger.Errorln(err)
		return err
	}

	for _, sub := range subs {
		// search the subscribed tweeter from DB by twitter user id
		// and convert to username that is accecptable by twitterscraper
		tu, err := tl.repo.GetTweeterOfChatSub(context.Background(), int(sub.ID))
		if err != nil {
			tl.logger.Warnln(err)
			continue
		}

		_, lasttweet, err := ScrapLastTweet(tu.Username)
		if err != nil {
			tl.logger.WithField("subscribee username", tu.Username).Warnln(err)
			continue
		}

		// new tweet
		if twurl := lasttweet; twurl != sub.LastTweet {
			// publish the new tweet url to event listeners
			tl.logger.
				WithField("tweet url", twurl).
				WithField("tochat", chatid).
				Println("updated tweet")

			tl.evhub.Fire("newtweet", chatid, twurl)

			// update state
			tl.repo.UpdateLastTweet(context.Background(), chatid, tu.ID, twurl)
		}

	}

	return nil

}
