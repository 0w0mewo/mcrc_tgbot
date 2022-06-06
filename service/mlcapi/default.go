package mlcapi

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

var apilist = map[string]any{
	API_BAN:       Ban,
	API_UNBAN:     UnBan,
	API_DELTEUSER: Delete,
	API_GETINFO:   GetInfo,
	API_CHANGEPW:  ChangePassword,
	API_RESET:     ResetHWID,
	API_REGISTER:  Register,
}

var defaultClient *MlcApiClient

func init() {
	defaultClient = NewMlcApiClient(os.Getenv("STAFF_TOKEN"), os.Getenv("MAN_TOKEN"))
}

// register user
func Register(username string, password string) error {
	return defaultClient.Register(username, password)

}

// ban user with reason
func Ban(username string, reason string) error {
	return defaultClient.Ban(username, reason)

}

// unban user
func UnBan(username string) error {
	return defaultClient.UnBan(username)

}

// delete user
func Delete(username string) error {
	return defaultClient.Delete(username)

}

// reset hwid
func ResetHWID(username string) error {
	return defaultClient.ResetHWID(username)

}

// change password
func ChangePassword(username string, newpassword string) error {
	return defaultClient.ChangePassword(username, newpassword)

}

// get info of user
func GetInfo(username string) (*UserInfo, error) {
	return defaultClient.GetInfo(username)

}

// call mlc api with api name
func ApiCall(apiname string, args ...string) (ret any, err error) {
	defer func() {
		if e := recover(); e != nil {
			ret = nil
			switch e := e.(type) {
			case error:
				err = e
			case string:
				_, e, _ = strings.Cut(e, ":")
				err = errors.New(e)
			}
		}
	}()
	apicall := reflect.ValueOf(apilist[apiname])
	if apicall.IsNil() {
		return nil, errors.New("unsupported API")
	}

	parmas := make([]reflect.Value, 0, len(args))

	for i := 0; i < reflect.TypeOf(apilist[apiname]).NumIn() && i < len(args); i++ {
		parmas = append(parmas, reflect.ValueOf(args[i]))

	}

	rets := apicall.Call(parmas)
	if len(rets) < 2 {
		if rets[0].Interface() == nil {
			return
		}
		return nil, rets[0].Interface().(error)
	}

	if rets[1].Interface() == nil {
		ret = rets[0].Interface()
		return
	}

	return nil, rets[1].Interface().(error)

}
