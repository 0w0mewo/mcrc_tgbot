package linesticker

import (
	"fmt"
)

const PNGURL = "http://dl.stickershop.line.naver.jp/stickershop/v1/sticker/%d/iphone/sticker@2x.png"
const GIFURL = "https://sdl-stickershop.line.naver.jp/products/0/0/1/%d/iphone/animation/%d@2x.png"
const METAINFOURL = "http://dl.stickershop.line.naver.jp/products/0/0/1/%d/android/productInfo.meta"

type PackMeta struct {
	PackageId int        `json:"packageId"`
	Stickers  []*Sticker `json:"stickers"`
	HasGif    bool       `json:"hasAnimation"`
}

type Sticker struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Id     int `json:"id"`
}

func (s *Sticker) Key(animated bool) string {
	var fileExt string
	if animated {
		fileExt = "gif"
	} else {
		fileExt = "png"
	}

	return fmt.Sprintf("%d.%s", s.Id, fileExt)
}
