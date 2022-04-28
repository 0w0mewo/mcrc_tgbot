package model

import "time"

// model of telegram messages
type TeleMsg struct {
	ChatName       string
	SenderUserName string
	Type           int
	Message        []byte
	TimeStamp      time.Time
}

type ChatTweetSubscription struct {
	Id        int // database record id
	ChatId    int64
	TweeterId string
	LastTweet string
}

type Chat struct {
	Id   int64
	Name string
}

type TweetUser struct {
	Id        string
	UserName  string
	LastTweet string
}
