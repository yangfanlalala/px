package wx

import (
	"encoding/base64"
	"encoding/json"
	"github.com/yangfanlalala/px/crypto"
)

type Watermark struct {
	AppID     string `json:"appid"`
	Timestamp int64  `json:"timestamp"`
}

func decrypt(cipherText, iv, session string, obj interface{}) error {
	cipherDecode, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return err
	}
	ivDecode, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return err
	}
	sessionDecode, err := base64.StdEncoding.DecodeString(session)
	if err != nil {
		return err
	}
	plainText, err := crypto.AesDecrypt(string(cipherDecode), string(sessionDecode), string(ivDecode))
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(plainText), obj); err != nil {
		return err
	}
	return nil
}
