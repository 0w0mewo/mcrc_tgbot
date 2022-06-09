package utils

// func DumpChatMessage(msgrepo persistent.StoredTeleMsgRepo, todir string, chatid int64) error {
// 	// total messages count
// 	cnt, err := msgrepo.Count(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	// paging
// 	pagesize := 1
// 	pagecnt := cnt / pagesize
// 	if pagecnt <= 0 {
// 		pagecnt = 1
// 	}

// 	// create text message storage
// 	fpath := filepath.Join(todir, fmt.Sprintf("%d.text_msg.txt", chatid))
// 	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0o644)
// 	if err != nil {
// 		return err
// 	}
// 	defer fd.Close()

// 	// go through every messages by chatid
// 	for cur := 0; cur < pagecnt; cur++ {
// 		msgs, err := msgrepo.GetChatMessages(context.Background(), chatid, cnt, pagesize)
// 		if err != nil {
// 			return err
// 		}

// 		for _, msg := range msgs {
// 			// text messages
// 			if msg.Type == persistent.MEDIA_TEXT {
// 				fmt.Fprintf(fd, "{%s} [%s]: \"%s\" at %s\n", msg.ChatName, msg.SenderUserName, msg.Message, msg.TimeStamp.Format(time.UnixDate))

// 			} else {
// 				fname := fmt.Sprintf("%s-%s-file-%d", msg.ChatName, msg.SenderUserName, msg.TimeStamp.Unix())
// 				err := ioutil.WriteFile(fname, msg.Message, 0o644)
// 				if err != nil {
// 					continue
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }
