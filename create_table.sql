CREATE TABLE IF NOT EXISTS `tweet_users`(`id` varchar(255) UNIQUE NOT NULL, `username` varchar(255) NOT NULL, PRIMARY KEY(`id`));
CREATE INDEX IF NOT EXISTS `tweetuser_username` ON `tweet_users`(`username`);
CREATE TABLE IF NOT EXISTS `chats`(`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL, `name` varchar(255) NOT NULL);
CREATE TABLE IF NOT EXISTS `messages`(`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL, `msg` blob NOT NULL, `type` integer NOT NULL, `timestamp` datetime NOT NULL, `chat_id` integer NULL, `sender_id` integer NULL, FOREIGN KEY(`chat_id`) REFERENCES `chats`(`id`) ON DELETE SET NULL, FOREIGN KEY(`sender_id`) REFERENCES `senders`(`id`) ON DELETE SET NULL);
CREATE TABLE IF NOT EXISTS `senders`(`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL, `username` varchar(255) NOT NULL);
CREATE INDEX IF NOT EXISTS `message_chat_id_sender_id` ON `messages`(`chat_id`, `sender_id`);
CREATE INDEX IF NOT EXISTS `message_chat_id` ON `messages`(`chat_id`);
CREATE TABLE IF NOT EXISTS `chat_tweet_subscriptions`(`id` integer PRIMARY KEY AUTOINCREMENT NOT NULL, `last_tweet` varchar(255) NOT NULL, `tweeter_id` varchar(255) NULL, `chat_id` integer NULL, FOREIGN KEY(`tweeter_id`) REFERENCES `tweet_users`(`id`) ON DELETE SET NULL, FOREIGN KEY(`chat_id`) REFERENCES `chats`(`id`) ON DELETE SET NULL);
CREATE INDEX IF NOT EXISTS `chattweetsubscription_chat_id_tweeter_id` ON `chat_tweet_subscriptions`(`chat_id`, `tweeter_id`);
