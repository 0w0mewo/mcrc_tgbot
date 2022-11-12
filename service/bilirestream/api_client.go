package bilirestream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/0w0mewo/mcrc_tgbot/utils"
)

type ApiClient struct {
	hc      *http.Client
	apilist map[string]any
}

func NewClient() *ApiClient {
	client := &ApiClient{
		hc: http.DefaultClient,
	}

	client.apilist = map[string]any{
		API_BILIBILI_RESTREAM_START:  "restream",
		API_BILIBILI_RESTREAM_STATUS: "status",
		API_BILIBILI_RESTREAM_STOP:   "stop",
	}

	return client
}

func (mr *ApiClient) BiliRestream(reqtype string, username string, token string, roomid string) (resp *RestreamResp, err error) {
	var url string

	switch reqtype {
	case API_BILIBILI_RESTREAM_STATUS:
		url = fmt.Sprintf("%s/status?user=%s", restreamEndpoint, username)
	case API_BILIBILI_RESTREAM_STOP:
		url = fmt.Sprintf("%s/stop?user=%s&token=%s", restreamEndpoint, username, token)
	case API_BILIBILI_RESTREAM_START:
		url = fmt.Sprintf("%s/restream?user=%s&token=%s&biliroom=%s", restreamEndpoint, username, token, roomid)
	default:
		return nil, ErrBadRequest
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp = &RestreamResp{}
	err = utils.HttpGetJson(ctx, mr.hc, url, resp)
	if err != nil {
		return
	}

	return

}

func (mr *ApiClient) Gentoken(username string) (resp *TokenResp, err error) {
	url := fmt.Sprintf("%s/%s", gentokenEndpoint, username)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp = new(TokenResp)

	err = utils.HttpDoWithProcessor(ctx, mr.hc, http.MethodPut, url, nil, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(resp)
	})

	return

}
