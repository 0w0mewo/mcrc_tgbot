package mlcapi

import "github.com/pkg/errors"

var ErrUnknown = errors.New("unknown error")
var ErrBadRequest = errors.New("invalid request type")
