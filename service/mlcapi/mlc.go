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
	apilist      map[string]any
}

func NewMlcApiClient(staffToken, managerToken string) *MlcApiClient {
	client := &MlcApiClient{
		staffToken:   staffToken,
		managerToken: managerToken,
		hc:           http.DefaultClient,
	}

	client.apilist = map[string]any{
		API_BAN:       client.Ban,
		API_UNBAN:     client.UnBan,
		API_DELTEUSER: client.Delete,
		API_GETINFO:   client.GetInfo,
		API_CHANGEPW:  client.ChangePassword,
		API_RESET:     client.ResetHWID,
		API_REGISTER:  client.Register,
		API_SETRANK:   client.SetRank,
	}

	return client
}

func (mr *MlcApiClient) SetStaffToken(token string) {
	mr.staffToken = token
}

func (mr *MlcApiClient) SetManagerToken(token string) {
	mr.managerToken = token
}

func (mr *MlcApiClient) buildStaffReq(httpmethod string, reqtype string, endpoint string, apireq UserGeneralReq) (req *http.Request, err error) {
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
		form.Add("rank", apireq.Rank)
	case API_SETRANK:
		form.Add("rank", apireq.Rank)
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

func (mr *MlcApiClient) do(httpmethod, reqtype string, endpoint string, r UserGeneralReq) (*UserGeneralResp, error) {
	req, err := mr.buildStaffReq(httpmethod, reqtype, endpoint, r)
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
func (mr *MlcApiClient) Register(username string, password string, rank string) (err error) {
	_, err = mr.do(http.MethodPost, API_REGISTER, staffEndpoint, UserGeneralReq{
		UserName: username,
		Password: password,
		Rank:     rank,
	})

	return
}

// ban user with reason
func (mr *MlcApiClient) Ban(username string, reason string) (err error) {
	_, err = mr.do(http.MethodPost, API_BAN, staffEndpoint, UserGeneralReq{
		UserName:    username,
		BanedReason: reason,
	})

	return
}

// unban user
func (mr *MlcApiClient) UnBan(username string) (err error) {
	_, err = mr.do(http.MethodPost, API_UNBAN, staffEndpoint, UserGeneralReq{
		UserName: username,
	})

	return
}

// delete user
func (mr *MlcApiClient) Delete(username string) (err error) {
	_, err = mr.do(http.MethodPost, API_DELTEUSER, staffEndpoint, UserGeneralReq{
		UserName: username,
	})

	return
}

// reset hwid
func (mr *MlcApiClient) ResetHWID(username string) (err error) {
	_, err = mr.do(http.MethodPost, API_RESET, staffEndpoint, UserGeneralReq{
		UserName: username,
	})

	return
}

// change password
func (mr *MlcApiClient) ChangePassword(username string, newpassword string) (err error) {
	_, err = mr.do(http.MethodPost, API_CHANGEPW, staffEndpoint, UserGeneralReq{
		UserName: username,
		Password: newpassword,
	})

	return
}

// get info of user
func (mr *MlcApiClient) GetInfo(username string) (*UserInfo, error) {
	resp, err := mr.do(http.MethodPost, API_GETINFO, staffEndpoint, UserGeneralReq{
		UserName: username,
	})
	if err != nil {
		return nil, err
	}

	return resp.Info, nil
}

// set rank of user
func (mr *MlcApiClient) SetRank(username string, rank string) error {
	_, err := mr.do(http.MethodPost, API_SETRANK, staffEndpoint, UserGeneralReq{
		UserName: username,
		Rank:     rank,
	})

	return err
}
