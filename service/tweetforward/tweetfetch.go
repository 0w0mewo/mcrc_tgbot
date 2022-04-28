package tweetforward

import (
	"errors"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/model"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

var (
	ErrSubscribed     = errors.New("tweeter is already subscribed")
	ErrNotSubscribed  = errors.New("tweeter is not subscribed")
	ErrNotExist       = errors.New("no such user")
	ErrTwitterScraper = errors.New("twitter API unkonwn error: tweet result is null")
)

func ScrapLastTweet(uname string) (twu *model.TweetUser, err error) {
	scrapLastTweet := func(uname string) (twu *model.TweetUser, err error) {
		tweets, _, err := twitterscraper.New().FetchTweets(uname, 1, "")
		if err != nil {
			return
		}

		if len(tweets) == 0 {
			err = ErrTwitterScraper
			return
		}

		lastTweetInfo := tweets[0]
		if lastTweetInfo == nil {
			err = ErrTwitterScraper
			return
		}
		twuid, twuname, lastTweet := lastTweetInfo.UserID, lastTweetInfo.Username, lastTweetInfo.PermanentURL

		if twuid == "" || twuname == "" {
			err = ErrNotExist
			return
		}

		if lastTweet == "" {
			err = ErrTwitterScraper
			return
		}

		twu = new(model.TweetUser)
		twu.Id = twuid
		twu.UserName = twuname
		twu.LastTweet = lastTweet

		return
	}

	// retry twice and delay 2 seconds for each trial if fail
	for retries := 0; retries < 2; retries++ {
		twu, err = scrapLastTweet(uname)
		if err != nil && !errors.Is(err, ErrNotExist) {
			time.Sleep(2 * time.Second)
			continue
		} else {
			break
		}
	}

	return
}
