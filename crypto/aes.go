package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
)

func AesEncrypt(plaintext, key, iv string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	plaintext = pkcs7padding(plaintext, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, []byte(iv))
	ciphertext := make([]byte, len(plaintext))
	blockModel.CryptBlocks(ciphertext, []byte(plaintext))
	return string(ciphertext), nil
}

func AesDecrypt(ciphertext, key, iv string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originText := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(originText, []byte(ciphertext))
	return pkcs7unpadding(string(originText)), nil
}

func AesQuickEncrypt(plaintext, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	iv := make([]byte, block.BlockSize())
	_, _ = rand.Read(iv)
	ciphertext, _ := AesEncrypt(plaintext, key, string(iv))
	dataMap := map[string]string{"iv": string(iv), "value": ciphertext}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func AesQuickDecrypt(cipher, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(cipher))
	if err != nil {
		return "", errors.New("BASE64 解码错误")
	}
	x := make(map[string]string)
	_ = json.Unmarshal(ciphertext, &x)
	iv, ok := x["iv"]
	if !ok {
		return "", errors.New("无法确定iv")
	}
	value, _ := x["value"]
	ivDecode, _ := base64.StdEncoding.DecodeString(iv)
	vaDecode, _ := base64.StdEncoding.DecodeString(value)
	return AesDecrypt(string(vaDecode), key, string(ivDecode))
}

func pkcs7padding(ciphertext string, blockSize int) string {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return ciphertext + string(padText)
}

func pkcs7unpadding(plaintext string) string {
	length := len(plaintext)
	padding := int(plaintext[length-1])
	return plaintext[:(length - padding)]
}
