package mlcapi

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

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
func ApiCall(client *MlcApiClient, apiname string, args ...string) (ret any, err error) {
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
	apicall := reflect.ValueOf(client.apilist[apiname])
	if apicall.IsNil() {
		return nil, errors.New("unsupported API")
	}

	parmas := make([]reflect.Value, 0, len(args))

	for i := 0; i < reflect.TypeOf(client.apilist[apiname]).NumIn() && i < len(args); i++ {
		parmas = append(parmas, reflect.ValueOf(args[i]))

	}

	rets := apicall.Call(parmas)

	if len(rets) < 2 {
		err = reflectValAsError(rets[0])
		if err != nil {
			ret = nil
			return
		}

		return
	}

	err = reflectValAsError(rets[1])
	if err != nil {
		ret = nil
		return
	}

	ret = rets[0].Interface()
	return
}

func reflectValAsError(val reflect.Value) error {
	v := val.Interface()
	switch err := v.(type) {
	case error:
		return err
	case nil:
		return nil
	default:
		return ErrUnknown
	}
}
