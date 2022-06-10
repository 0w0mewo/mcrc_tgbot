package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/persistent"
)

func DumpChatMessage(msgrepo persistent.StoredTeleMsgRepo, todir string, chatid int64) error {
	// total messages count
	cnt, err := msgrepo.Count(context.Background())
	if err != nil {
		return err
	}

	// create text message storage
	fpath := filepath.Join(todir, fmt.Sprintf("%d.text_msg.txt", chatid))
	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer fd.Close()

	// paging
	pager := NewPager(cnt, 10)

	// go through every messages by chatid
	pager.Iter(func(p *Pager) error {
		msgs, err := msgrepo.GetChatMessages(context.Background(), chatid, p.Page, 10)
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			// text messages
			if msg.Type == persistent.MEDIA_TEXT {
				fmt.Fprintf(fd, "{%d} [%d]: \"%s\" at %s\n", msg.ChatID.Int64, msg.SenderID.Int64, msg.MSG, msg.Timestamp.Format(time.UnixDate))

			} else {
				fname := fmt.Sprintf("%d-%d-file-%d", msg.ChatID.Int64, msg.SenderID.Int64, msg.Timestamp.Unix())
				err := ioutil.WriteFile(fname, msg.MSG, 0o644)
				if err != nil {
					continue
				}
			}
		}

		return err
	})

	return nil
}

type Pager struct {
	pages    int // total page
	pagesize int // elements per page

	Page int // current page number
	Cur  int // current offset
}

func NewPager(nelements int, pagesize int) *Pager {
	pages := nelements / pagesize
	if pages <= 0 {
		pages = 1
	}
	return &Pager{
		pages:    pages,
		pagesize: pagesize,
		Cur:      0,
		Page:     1,
	}
}

// iterate pages
func (p *Pager) Iter(fn func(p *Pager) error) (err error) {
	for cur, page := 0, 1; p.Page < p.pages; page++ {
		cur = p.pagesize * (page - 1)

		p.Cur = cur
		p.Page = page
		err = fn(p)
	}

	return
}
