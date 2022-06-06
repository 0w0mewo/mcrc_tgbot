package mlcapi

type GeneralResp struct {
	Status string `json:"status"`
	Msg    string `json:"error_msg,omitempty"`
}

// general purpose user management request
type UserGeneralReq struct {
	UserName    string
	Password    string
	BanedReason string
}

// user info response
type UserGeneralResp struct {
	GeneralResp
	Info *UserInfo `json:"info,omitempty"`
}

// user info
type UserInfo struct {
	UserName          string `json:"username"`
	BanedReason       string `json:"banReason"`
	Hwid              string `json:"hardwareId"`
	LastIP            string `json:"lastLoginIP"`
	State             string `json:"userState"`
	PasswordPlaintext string `json:"password"`
	LastLoginTime     int64  `json:"lastLoginTime"`
}
