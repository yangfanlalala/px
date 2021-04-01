package wx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	WeChatURLToken                = "https://api.weixin.qq.com/cgi-bin/token"
	WeChatURLACodeUnlimited       = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
	WeChatURLSubscribeMessageSend = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"
	WeChatURLCodeToSession        = "https://api.weixin.qq.com/sns/jscode2session"

	WeChatSuccessCode = 0

	MiniProgramStateDeveloper = "developer"
	MiniProgramStateTrial     = "trial"
	MiniProgramStateFormal    = "formal"
)

var (
	HttpClientIsNil = errors.New("http client is nil")
)

type AccessToken struct {
	Token   string
	Expired int64
}

type MiniProgramClient struct {
	AppID      string
	AppSecret  string
	httpClient *http.Client
}

func NewMiniProgramClient(ak, as string, cli *http.Client) *MiniProgramClient {
	return &MiniProgramClient{
		AppID:      ak,
		AppSecret:  as,
		httpClient: cli,
	}
}

//获取服务器访问令牌
func (wx *MiniProgramClient) GetAccessToken() (*AccessToken, error) {
	if wx.httpClient == nil {
		return nil, HttpClientIsNil
	}
	val := &url.Values{}
	val.Set("grant_type", "client_credential")
	val.Set("appid", wx.AppID)
	val.Set("secret", wx.AppSecret)
	req, err := http.NewRequest(http.MethodGet, WeChatURLToken+"?"+val.Encode(), nil)
	if err != nil {
		return nil, err
	}
	res, err := wx.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r := &struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		ErrorCode    int64  `json:"errcode"`
		ErrorMessage string `json:"errmsg"`
	}{}
	if err = json.Unmarshal(body, r); err != nil {
		return nil, err
	}
	if r.ErrorCode != 0 {
		return nil, fmt.Errorf("request weixin service failed code[%d] messge[%s]", r.ErrorCode, r.ErrorMessage)
	}
	return &AccessToken{
		Token:   r.AccessToken,
		Expired: r.ExpiresIn,
	}, nil
}

//获取小程序码
//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func (wx *MiniProgramClient) GetWxACodeUnlimited(ac string, scene, page string, width uint32) (io.Closer, error) {
	if wx.httpClient == nil {
		return nil, HttpClientIsNil
	}
	params := struct {
		Scene string `json:"scene"`
		Page  string `json:"page"`
		Width uint32 `json:"width"`
	}{Scene: scene, Page: page, Width: width}
	content, _ := json.Marshal(params)
	req, err := http.NewRequest(http.MethodPost, WeChatURLACodeUnlimited+"?access_token="+ac, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	res, err := wx.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "image") {
		return res.Body, nil
	}
	body, err := ioutil.ReadAll(res.Body)
	r := &struct {
		ErrorCode    int32  `json:"errcode"`
		ErrorMessage string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("request weixin service failed code[%d] message[%s]", r.ErrorCode, r.ErrorMessage)
}

type TemplateMessage struct {
	TemplateID string      `json:"template_id"`
	Page       string      `json:"page"`
	Data       interface{} `json:"data"`
	State      string      `json:"miniprogram_state"`
}

//发送订阅消息
func (wx *MiniProgramClient) SendSubscribeMessage(ac string, msg TemplateMessage) error {
	if wx.httpClient == nil {
		return HttpClientIsNil
	}
	content, _ := json.Marshal(msg)
	request, err := http.NewRequest(http.MethodPost, WeChatURLSubscribeMessageSend+"?access_token="+ac, bytes.NewReader(content))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	res, err := wx.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed, code[%d], msg[%s]", res.StatusCode, body)
	}
	result := &struct {
		ErrorCode    int32  `json:"errCode"`
		ErrorMessage string `json:"errMsg"`
	}{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.ErrorCode != 0 {
		return fmt.Errorf("wx service return error, error code[%d], error message[%s]", result.ErrorCode, result.ErrorMessage)
	}
	return nil
}

type MiniProgramSession struct {
	OpenID     string
	SessionKey string
}

//获取Session
func (wx *MiniProgramClient) GetSession(code string) (*MiniProgramSession, error) {
	if wx.httpClient == nil {
		return nil, HttpClientIsNil
	}
	url := WeChatURLCodeToSession + "?js_code=" + code + "&appid=" + wx.AppID + "&secret=" + wx.AppSecret + "&grant_type=authorization_code"
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := wx.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request wx service failed, status code[%d], body[%d]", response.StatusCode, body)
	}
	reply := &struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int64  `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}{}
	if err = json.Unmarshal(body, reply); err != nil {
		return nil, err
	}
	if reply.ErrCode != WeChatSuccessCode {
		return nil, fmt.Errorf("request wx service failed, error code[%d], error message[%s]", reply.ErrCode, reply.ErrMsg)
	}
	return &MiniProgramSession{OpenID: reply.OpenID, SessionKey: reply.SessionKey}, nil
}

type MiniUserInformation struct {
	Nickname  string    `json:"nickName"`
	AvatarURL string    `json:"avatarUrl"`
	Gender    uint8     `json:"gender"`
	Country   string    `json:"country"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Language  string    `json:"language"`
	Watermark Watermark `json:"watermark"`
}

//解析用户信息
func (wx *MiniProgramClient) GetUserInformation(cipher, iv, key string) (*MiniUserInformation, error) {
	info := &MiniUserInformation{}
	if err := decrypt(cipher, iv, key, info); err != nil {
		return nil, err
	}
	if wx.AppID != info.Watermark.AppID {
		return nil, fmt.Errorf("appid not right")
	}
	return info, nil
}

type MiniUserPhoneInformation struct {
	PhoneNumber     string    `json:"phoneNumber"`
	PurePhoneNumber string    `json:"purePhoneNumber"`
	CountryCode     string    `json:"countryCode"`
	Watermark       Watermark `json:"watermark"`
}

//解析手机号码
func (wx *MiniProgramClient) GetPhoneNumber(cipher, iv, key string) (*MiniUserPhoneInformation, error) {
	phone := &MiniUserPhoneInformation{}
	if err := decrypt(cipher, iv, key, phone); err != nil {
		return nil, err
	}
	if wx.AppID != phone.Watermark.AppID {
		return nil, fmt.Errorf("appid not right")
	}
	return phone, nil
}
