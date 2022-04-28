package tweetforward

import (
	"context"
	"errors"
	"sync"
	"time"

	constant "github.com/0w0mewo/mcrc_tgbot/const"
	"github.com/0w0mewo/mcrc_tgbot/model"
	"github.com/0w0mewo/mcrc_tgbot/persistent"
	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/agoalofalife/event"
	"github.com/sirupsen/logrus"
)

var defaultForwarder *tweetForwarder

func init() {
	defaultForwarder = newTweetForwarder()

}

func Subscribe(chat model.Chat, username string) error {
	return defaultForwarder.Subscribe(chat, username)
}

func Unsubscribe(chatid int64, username string) error {
	return defaultForwarder.Unsubscribe(chatid, username)
}

func GetSubscriptions(chatid int64) ([]*model.TweetUser, error) {
	return defaultForwarder.GetSubscriptions(chatid)
}

func DoWhenSet(cb any) {
	defaultForwarder.DoWhenSet(cb)
}

func Shutdown() {

	defaultForwarder.Shutdown()
}

type tweetForwarder struct {
	repo   persistent.ChatTweetSubRepo
	logger *logrus.Entry
	quit   chan bool
	wg     *sync.WaitGroup
	evhub  *event.Dispatcher
}

func newTweetForwarder() *tweetForwarder {

	ret := &tweetForwarder{
		repo:   persistent.NewChatTweetSubSqlStorage(persistent.DefaultDBConn),
		logger: utils.NewLogger().WithField("service", "tweetlistener"),
		wg:     &sync.WaitGroup{},
		evhub:  event.New(),
		quit:   make(chan bool),
	}

	ret.wg.Add(1)
	go ret.watch(constant.TWEET_REFRESH_INTERVAL)

	return ret
}

func (tl *tweetForwarder) GetSubscriptions(chatid int64) ([]*model.TweetUser, error) {
	return tl.repo.GetAllSubscribeeByChatId(context.Background(), chatid)
}

func (tl *tweetForwarder) Subscribe(chat model.Chat, tweeter string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DB_CRUD_TIMEOUT)
	defer cancel()

	twu, err := ScrapLastTweet(tweeter)
	if err != nil {
		return err
	}

	err = tl.repo.Create(ctx, chat, *twu)
	if errors.Is(err, persistent.ErrExist) {
		return ErrSubscribed
	}

	return err

}

func (tl *tweetForwarder) Unsubscribe(chatid int64, tweeter string) error {
	ctx, cancel := context.WithTimeout(context.Background(), constant.DB_CRUD_TIMEOUT)
	defer cancel()

	twu, err := ScrapLastTweet(tweeter)
	if err != nil {
		return err
	}

	err = tl.repo.Remove(ctx, chatid, twu.Id)
	if errors.Is(err, persistent.ErrNotExist) {
		return ErrNotSubscribed
	}

	return err
}

func (tl *tweetForwarder) DoWhenSet(cb any) {
	tl.evhub.Add("newtweet", cb)
}

func (tl *tweetForwarder) Shutdown() {
	tl.quit <- true

	err := tl.repo.Close()
	if err != nil {
		tl.logger.Error(err)
	}
	tl.wg.Wait()
}

func (tl *tweetForwarder) watch(interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer func() {
		ticker.Stop()
		tl.wg.Done()
		close(tl.quit)
	}()

	for {
		select {
		case <-tl.quit:
			return

		case <-ticker.C:
			var wg sync.WaitGroup // for chats polling

			// get all chats which subscribed various tweeters
			chats, err := tl.repo.GetAllChat(context.Background())
			if err != nil {
				tl.logger.Errorln(err)
				continue
			}

			// go through each chat
			for _, chatid := range chats {
				wg.Add(1)
				go tl.updateChatSubs(chatid, &wg)
			}

			wg.Wait()
		}

	}
}

func (tl *tweetForwarder) updateChatSubs(chatid int64, joinwg *sync.WaitGroup) error {
	defer joinwg.Done()

	subs, err := tl.repo.GetAllChatSub(context.Background(), chatid)
	if err != nil {
		tl.logger.Errorln(err)
		return err
	}

	// go through each subcribee in a chat
	for _, sub := range subs {
		// search the subscribed tweeter from DB by twitter user id
		// and convert to username that is accecptable by twitterscraper
		tu, err := tl.repo.GetTweeterOfChatSub(context.Background(), sub.Id)
		if err != nil {
			tl.logger.Warnln(err)
			continue
		}
		twu, err := ScrapLastTweet(tu.UserName)
		if err != nil {
			tl.logger.WithField("subscribee username", tu.UserName).Warnln(err)
			continue
		}

		// new tweet
		if twurl := twu.LastTweet; twurl != sub.LastTweet {
			// publish the new tweet url to event listeners
			tl.logger.
				WithField("tweet url", twurl).
				WithField("tochat", chatid).
				Println("updated tweet")

			tl.evhub.Fire("newtweet", chatid, twurl)

			// update state
			tl.repo.UpdateLastTweet(context.Background(), chatid, tu.Id, twurl)
		}

	}

	return nil

}
