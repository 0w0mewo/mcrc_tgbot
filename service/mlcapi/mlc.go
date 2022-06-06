package mlcapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type MlcApiClient struct {
	staffToken   string
	managerToken string
	hc           *http.Client
}

func NewMlcApiClient(staffToken, managerToken string) *MlcApiClient {
	return &MlcApiClient{
		staffToken:   staffToken,
		managerToken: managerToken,
		hc:           http.DefaultClient,
	}
}

func (mr *MlcApiClient) buildReq(httpmethod string, reqtype string, endpoint string, apireq UserGeneralReq) (req *http.Request, err error) {
	form := url.Values{}
	form.Add("type", reqtype)
	form.Add("username", apireq.UserName)

	switch reqtype {
	case API_BAN:
		form.Add("reason", apireq.BanedReason)
	case API_CHANGEPW:
		form.Add("new_password", apireq.Password)
	case API_REGISTER:
		form.Add("password", apireq.Password)
	case API_RESET:
	case API_UNBAN:
	case API_GETINFO:
	case API_DELTEUSER:
	default:
		return nil, ErrBadRequest
	}

	req, err = http.NewRequest(httpmethod, endpoint, strings.NewReader(form.Encode()))

	// set token header
	if strings.Contains(req.URL.Path, "manager") {
		req.Header.Add("token", mr.managerToken)
	} else {
		req.Header.Add("token", mr.staffToken)
	}

	req.Header.Add("User-Agent", "NMSL")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return

}

func (mr *MlcApiClient) do(reqtype string, endpoint string, r UserGeneralReq) (*UserGeneralResp, error) {
	req, err := mr.buildReq("POST", reqtype, endpoint, r)
	if err != nil {
		return nil, errors.Wrap(err, "buildreq")
	}

	resp, err := mr.hc.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do")
	}
	defer resp.Body.Close()

	var res UserGeneralResp

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "ioread")
	}

	// response content-type is NOT "application/json"
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, errors.Wrap(err, "json")
	}

	switch status := strings.ToLower(res.Status); status {
	case "success":
		return &res, nil
	default:
		return nil, errors.New(res.Msg)
	}
}

// register user
func (mr *MlcApiClient) Register(username string, password string) (err error) {
	_, err = mr.do(API_REGISTER, staffEndpoint, UserGeneralReq{
		UserName: username,
		Password: password,
	})

	return
}

// ban user with reason
func (mr *MlcApiClient) Ban(username string, reason string) (err error) {
	_, err = mr.do(API_BAN, staffEndpoint, UserGeneralReq{
		UserName:    username,
		BanedReason: reason,
	})

	return
}

// unban user
func (mr *MlcApiClient) UnBan(username string) (err error) {
	_, err = mr.do(API_UNBAN, staffEndpoint, UserGeneralReq{
		UserName: username,
	})

	return
}

// delete user
func (mr *MlcApiClient) Delete(username string) (err error) {
	_, err = mr.do(API_DELTEUSER, staffEndpoint, UserGeneralReq{
		UserName: username,
	})

	return
}

// reset hwid
func (mr *MlcApiClient) ResetHWID(username string) (err error) {
	_, err = mr.do(API_RESET, staffEndpoint, UserGeneralReq{
		UserName: username,
	})

	return
}

// change password
func (mr *MlcApiClient) ChangePassword(username string, newpassword string) (err error) {
	_, err = mr.do(API_CHANGEPW, staffEndpoint, UserGeneralReq{
		UserName: username,
		Password: newpassword,
	})

	return
}

// get info of user
func (mr *MlcApiClient) GetInfo(username string) (*UserInfo, error) {
	resp, err := mr.do(API_GETINFO, staffEndpoint, UserGeneralReq{
		UserName: username,
	})
	if err != nil {
		return nil, err
	}

	return resp.Info, nil
}
