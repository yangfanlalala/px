package wx

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	WeChatURLWxACode        = "https://api.weixin.qq.com/wxa/getwxacode?access_token=ACCESS_TOKEN"
	WeChatURLWxACodeUnlimit = "https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=ACCESS_TOKEN"
)

type accessToken struct {
	token   string
	expired int64
}

type MiniProgram struct {
	appID      string
	appSecret  string
	token      accessToken
	httpClient *http.Client
}

func (wx *MiniProgram) GetAccessToken() (string, error) {
	return "", nil
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
func (wx *MiniProgram) GetWxACodeUnlimit(scene, page string, width uint32) {
	params := struct {
		Scene string `json:"scene"`
		Page  string `json:"page"`
		Width uint32 `json:"width"`
	}{Scene: scene, Page: page, Width: width}
	content, _ := json.Marshal(params)
	req, err := http.NewRequest(http.MethodPost, WeChatURLWxACodeUnlimit, bytes.NewReader(content))
	if err != nil {
		return
	}
	res, err := wx.httpClient.Do(req)
	if err != nil {
		return
	}
	defer func() { _ = res.Body.Close() }()
}

func (wx *MiniProgram) SendSubscribeMessage() error {
	return nil
}
