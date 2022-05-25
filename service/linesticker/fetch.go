package linesticker

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/0w0mewo/mcrc_tgbot/utils"
	"github.com/sirupsen/logrus"
)

// sticker package fetcher
type Fetcher struct {
	client *http.Client
	logger *logrus.Entry
}

// new sticker package fetcher
func NewFetcher(ctx context.Context, client *http.Client) *Fetcher {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	ret := &Fetcher{
		client: client,
		logger: logger.WithField("service", "line-sticker-fetcher"),
	}

	return ret
}

// fetch stickers package and save it to saveToDir
func (wk *Fetcher) SaveStickers(ctx context.Context, packid int, qqTrans bool) ([]byte, error) {
	// get stickers pack meta
	var pack PackMeta
	url := fmt.Sprintf(METAINFOURL, packid)

	err := utils.HttpGetJson(ctx, wk.client, url, &pack)
	if err != nil {
		return nil, err
	}

	animated := pack.HasGif

	zipres := bytes.NewBuffer(nil)
	zipper := zip.NewWriter(zipres)

	var mu sync.Mutex

	// process downloaded sticker
	stickerStorer := func(r io.Reader, s *Sticker) error {
		var folderName string

		if animated {
			folderName = "animated"
		} else {
			folderName = "not-animated"
		}

		// path: ./<packid>/<animated | not-animated>/<sticker>.<png | gif>
		stickerFolder := filepath.Join(".", folderName)
		path := filepath.Join(strconv.Itoa(packid), stickerFolder, s.Key(animated))

		mu.Lock()
		defer mu.Unlock()

		zfd, err := zipper.Create(path)
		if err != nil {
			return err
		}

		if animated || !qqTrans {
			_, err = io.Copy(zfd, r)
			return err
		} else {
			// support for transparency when import to qq
			// qq only recognises transparency background while the image format is gif
			return utils.PngToGif(zfd, r)
		}

	}
	var downloaders sync.WaitGroup

	// fetch and save stickers pack
	for _, s := range pack.Stickers {
		downloaders.Add(1)
		go func(s *Sticker) error {
			defer downloaders.Done()

			sid := s.Id
			err := fetchSticker(wk.client, packid, sid, animated, func(r io.Reader) error {
				return stickerStorer(r, s)
			})
			if err != nil {
				return err
			}

			return nil
		}(s)

	}

	downloaders.Wait()

	zipper.Close()

	return zipres.Bytes(), nil

}

func fetchSticker(client *http.Client, packid int, stickerId int, isAnimated bool, fn func(r io.Reader) error) error {
	var stickerUrl string

	if isAnimated {
		stickerUrl = fmt.Sprintf(GIFURL, packid, stickerId)
	} else {
		stickerUrl = fmt.Sprintf(PNGURL, stickerId)
	}

	return utils.HttpGetWithProcessor(context.Background(), client, stickerUrl, fn)

}
