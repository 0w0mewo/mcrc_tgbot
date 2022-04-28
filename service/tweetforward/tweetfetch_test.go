package tweetforward_test

import (
	"testing"

	"github.com/0w0mewo/mcrc_tgbot/service/tweetforward"
)

func TestScrapTweet(t *testing.T) {
	testUserName := "Johnny_Lindecis"
	twu, err := tweetforward.ScrapLastTweet(testUserName)
	if err != nil {
		t.Fatal(err)
	}

	if twu.UserName != testUserName {
		t.Fatal()
	}

	t.Logf("id = %s, uname = %s, last tweet = %s", twu.Id, twu.UserName, twu.LastTweet)
}
