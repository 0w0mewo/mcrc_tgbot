package bilirestream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	hresp, err := mr.hc.Do(req)
	if err != nil {
		return
	}
	defer hresp.Body.Close()

	var rresp RestreamResp

	err = json.NewDecoder(hresp.Body).Decode(&rresp)
	if err != nil {
		return
	}

	resp = &rresp
	return

}

func (mr *ApiClient) Gentoken(username string) (resp *TokenResp, err error) {
	url := fmt.Sprintf("%s/%s", gentokenEndpoint, username)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return
	}

	hresp, err := mr.hc.Do(req)
	if err != nil {
		return
	}
	defer hresp.Body.Close()

	var rresp TokenResp
	err = json.NewDecoder(hresp.Body).Decode(&rresp)
	if err != nil {
		return
	}

	resp = &rresp
	return

}
