package wx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
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

type miniCodeToSession struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int64  `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (wx *miniProgram) GetSession(code string) (openID, sessionKey string, err error) {
	fullURL := miniURLCodeToSession + "?js_code=" + code + "&appid=" + wx.appID + "&secret=" + wx.secret + "&grant_type=authorization_code"
	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", "", err
	}
	if wx.httpClient == nil {
		wx.httpClient = &http.Client{Timeout: 2}
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
	reply := &miniCodeToSession{}
	if err = json.Unmarshal(body, reply); err != nil {
		return "", "", err
	}
	if reply.ErrCode != miniSuccessCode {
		return "", "", fmt.Errorf("request we chat service failed, code[%d], message[%s]", reply.ErrCode, reply.ErrMsg)
	}
	return reply.OpenID, reply.SessionKey, nil
}

type MiniUserInfo struct {
	Nickname  string `json:"nickName"`
	AvatarURL string `json:"avatarUrl"`
	Gender    uint8  `json:"gender"`
	//Country string `json:"country"`
	//Province string `json:"province"`
	//City string `json:"city"`
	//Language string `json:"language"`
	Watermark Watermark `json:"watermark"`
}

func (wx *miniProgram) GetUserInfo(cipher, iv, key string) (*MiniUserInfo, error) {
	info := &MiniUserInfo{}
	if err := decrypt(cipher, iv, key, info); err != nil {
		return nil, err
	}
	if wx.appID != info.Watermark.AppID {
		return nil, fmt.Errorf("appid not right")
	}
	return info, nil
}

type MiniUserPhone struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

func (wx *miniProgram) GetPhone(cipher, iv, key string) (*MiniUserPhone, error) {
	phone := &MiniUserPhone{}
	if err := decrypt(cipher, iv, key, phone); err != nil {
		return nil, err
	}
	if wx.appID != phone.Watermark.AppID {
		return nil, fmt.Errorf("appid not right")
	}
	return phone, nil
}
