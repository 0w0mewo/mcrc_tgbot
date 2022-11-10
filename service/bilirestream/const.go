package bilirestream

import "errors"

const restreamEndpoint = "http://192.168.11.6:8888/bilirestream"
const gentokenEndpoint = "http://192.168.11.6:8888/token/gen"

var ErrUnknown = errors.New("unknown error")

const (
	API_GEN_TOKEN                = "gentoken"
	API_BILIBILI_RESTREAM_START  = "do"
	API_BILIBILI_RESTREAM_STOP   = "stop"
	API_BILIBILI_RESTREAM_STATUS = "status"
)
