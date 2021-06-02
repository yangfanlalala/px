package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

func RSAPKCS1Encrypt(plain, pk string) {

}

func RSAPKCS1Decrypt(cipher, pk string) {

}

func RSAPKCS8Encrypt(plain, pk string) {

}

func RSAPKCS8Decrypt(cipher, pk string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pk))
	if block == nil {
		return nil, errors.New("block is empty")
	}
	x, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("$2")
		return nil, err
	}
	ct, _ := base64.StdEncoding.DecodeString(cipher)
	rst, err := rsa.DecryptPKCS1v15(rand.Reader, x.(*rsa.PrivateKey), ct)
	if err != nil {
		fmt.Println("$1")
		return nil, err
	}
	return rst, nil
}
