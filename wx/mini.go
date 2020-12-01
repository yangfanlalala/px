package wx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	miniHost             = "https://api.weixin.qq.com"
	miniURLCodeToSession = "https://api.weixin.qq.com/sns/jscode2session"

	miniSuccessCode = 0
)

type miniProgram struct {
	appID      string
	secret     string
	httpClient *http.Client
}

func NewMiniProgramClient(cli *http.Client, appID, secret string) *miniProgram {
	return &miniProgram{
		appID:      appID,
		secret:     secret,
		httpClient: cli,
	}
}

type codeToSessionReply struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (wx *miniProgram) GetSession(code string) (openID, sessionKey string, err error) {
	fullURL := miniURLCodeToSession + "?js_code=" + code + "&appid=" + wx.appID + "&secret=" + wx.secret + "&grant_type=authorization_code"
	fmt.Println(fullURL)
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", "", err
	}
	resp, err := wx.httpClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	reply := &codeToSessionReply{}
	if err = json.Unmarshal(body, reply); err != nil {
		return "", "", err
	}
	if reply.ErrCode != miniSuccessCode {
		return "", "", fmt.Errorf("request we chat service failed, code[%d], message[%s]", reply.ErrCode, reply.ErrMsg)
	}
	return reply.OpenID, reply.SessionKey, nil
}
