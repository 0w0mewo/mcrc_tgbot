package bilirestream

import "github.com/pkg/errors"

// response for user creation/refresh token
type TokenResp struct {
	Code     int    `json:"code"`
	UserName string `json:"username"`
	Token    string `json:"token,omitempty"`
}

// srs auth resp
// expected user to stream to
// rtmp://xxxx/username/stream?token=TOKEN
type SrsReq struct {
	Action   string `json:"action"`
	UserName string `json:"app"`
	Stream   string `json:"stream"`
	Param    string `json:"param"`
}

type CommonResp struct {
	Code   int    `json:"code"`
	ErrMsg string `json:"err_msg"`
}

type RestreamResp struct {
	CommonResp
	Url string `json:"url"`
}

type BiliRestreamReq struct {
	User   string `query:"user" form:"user"`
	Token  string `query:"token" form:"token"`
	RoomId string `query:"biliroom" form:"biliroom"`
}

var ErrBadRequest = errors.New("invalid request type")
