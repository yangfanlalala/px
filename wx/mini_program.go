package wx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	WeChatURLToken = "https://api.weixin.qq.com/cgi-bin/token"
	WeChatURLWxACode        = "https://api.weixin.qq.com/wxa/getwxacode"
	WeChatURLWxACodeUnlimit = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
)

type accessToken struct {
	token   string
	expired int64
}

type lck struct {
	locked bool
	lck1 sync.Mutex
	lck2 sync.Mutex
}

func (l *lck) lock() bool {
	l.lck1.Lock()
	defer l.lck1.Unlock()
	if l.locked == false {
		l.locked = true
		l.lck2.Lock()
		return true
	}
	return false
}

func (l *lck) unlock() {
	l.lck1.Lock()
	defer l.lck1.Unlock()
	l.locked = false
	l.lck2.Unlock()
	return
}

type MiniProgram struct {
	appID      string
	appSecret  string
	token      accessToken
	httpClient *http.Client
	lock  lck
}

func NewMiniProgram (appID, appSecret string, cli *http.Client) *MiniProgram {
	return &MiniProgram{
		appID:      appID,
		appSecret:  appSecret,
		httpClient: cli,
	}
}

func (wx *MiniProgram) GetAccessToken() (string, error) {
	now := time.Now().Unix()
	//没问题，就直接返回
	if wx.token.expired > now {
		return wx.token.token, nil
	}
	if !wx.lock.lock() {
		time.Sleep(50 * time.Millisecond)
		return wx.GetAccessToken()
	}
	defer wx.lock.unlock()
	val := &url.Values{}
	val.Set("grant_type", "client_credential")
	val.Set("appid", wx.appID)
	val.Set("secret", wx.appSecret)
	req, err := http.NewRequest(http.MethodGet, WeChatURLToken + "?" + val.Encode(), nil)
	if err != nil {
		return "", err
	}
	res, err := wx.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {_ = res.Body.Close()}()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil
	}
	r := &struct {
		AccessToken string `json:"access_token"`
		ExpiresIn int64 `json:"expires_in"`
		ErrorCode int64 `json:"errcode"`
		ErrorMessage string `json:"errmsg"`
	}{}
	if err = json.Unmarshal(body, r); err != nil {
		return "", nil
	}
	if r.ErrorCode != 0 {
		return "", fmt.Errorf("request weixin service failed code[%d] messge[%s]", r.ErrorCode, r.ErrorMessage)
	}
	wx.token = accessToken{
		token:   r.AccessToken,
		expired: now + r.ExpiresIn - 1200,
	}
	return r.AccessToken, nil
}

func (wx *MiniProgram) GetWxACode(path string, width uint32) error {
	type s struct {
		Path  string `json:"path"`
		Width uint32 `json:"width"`
	}
	params := &s{Path: path, Width: width}
	content, _ := json.Marshal(params)
	req, err := http.NewRequest("POST", WeChatURLWxACode, bytes.NewReader(content))
	if err != nil {
		return err
	}
	res, err := wx.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = res.Body.Close() }()
	type r struct {
	}
	return nil
}

//https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/qr-code/wxacode.getUnlimited.html
func (wx *MiniProgram) GetWxACodeUnlimit(scene, page string, width uint32) (io.Closer, error) {
	params := struct {
		Scene string `json:"scene"`
		Page  string `json:"page"`
		Width uint32 `json:"width"`
	}{Scene: scene, Page: page, Width: width}
	content, _ := json.Marshal(params)
	ac, err := wx.GetAccessToken()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, WeChatURLWxACodeUnlimit + "?access_token=" + ac, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	res, err := wx.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix( contentType, "image") {
		return res.Body, nil
	}
	body, err := ioutil.ReadAll(res.Body)
	r := &struct {
		ErrorCode int32 `json:"errcode"`
		ErrorMessage string `json:"errmsg"`
	}{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("request weixin service failed code[%d] message[%s]", r.ErrorCode, r.ErrorMessage)
}

func (wx *MiniProgram) SendSubscribeMessage() error {
	return nil
}
