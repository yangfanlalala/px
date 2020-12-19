package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/yangfanlalala/px/crypto"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	OssDefaultExpireInSec = 5
)

type oss struct {
	ak       string
	as       string
	bucket   string
	endpoint string
	cli      *http.Client
}

type errResponse struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	RequestID string `xml:"RequestId"`
	HostID    string `xml:"HostId"`
}

func (e errResponse) Error() string {
	return "oss request error code[" + e.Code + "], message[" + e.Message + "], request_id[" + e.RequestID + "], host_id[" + e.HostID + "]"
}

func NewOss(ak, as, bucket, endpoint string, cli *http.Client) *oss {
	return &oss{
		ak:       ak,
		as:       as,
		bucket:   bucket,
		endpoint: endpoint,
		cli:      cli,
	}
}

func (o *oss) PutObject(data io.Reader, obj string) error {
	resp, err := o.do(http.MethodPut, obj, data)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		r := &errResponse{}
		e := xml.Unmarshal(body, r)
		if e != nil {
			return e
		}
		return r
	}
	fmt.Println("#23", string(body))
	return nil
}

func (o *oss) GetObjectURL(obj string, method string, expireInSec int64) string {
	if expireInSec <= 0 {
		expireInSec = OssDefaultExpireInSec
	}
	obj = strings.TrimLeft(obj, "/")
	//过期时间
	e := time.Now().Unix() + expireInSec
	es := strconv.FormatInt(e, 10)
	//签名字符串
	s := method + "\n\n\n" + es + "\n/" + o.bucket + "/" + obj
	signed := url.QueryEscape(base64.StdEncoding.EncodeToString(crypto.HmacSha1(s, o.as)))
	return "https://" + o.bucket + "." + o.endpoint + "/" + obj + "?Expires=" + es + "&OSSAccessKeyId=" + o.ak + "&Signature=" + signed
}

func (o *oss) GetImageObjectURL(obj string, expireInSec int64, options map[string]string) string {
	header := http.Header{}
	for k, v := range options {
		header.Add(k, v)
	}
	s := o.canonicalize(header)
	if expireInSec <= 0 {
		expireInSec = OssDefaultExpireInSec
	}
	obj = strings.TrimLeft(obj, "/")
	e := time.Now().Unix() + expireInSec
	es := strconv.FormatInt(e, 10)
	s = http.MethodGet + "\n\n\n" + es + "\n/" + s + o.bucket + "/" + obj
	signed := url.QueryEscape(base64.StdEncoding.EncodeToString(crypto.HmacSha1(s, o.as)))
	return "https://" + o.bucket + "." + o.endpoint + "/" + obj + "?Expires=" + es + "&OSSAccessKeyId=" + o.ak + "&Signature=" + signed
}

func (o *oss) DeleteObject(obj string) error {
	resp, err := o.do(http.MethodDelete, obj, nil)
	if err != nil {
		return err
	}
	defer func() { resp.Body.Close() }()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		r := &errResponse{}
		e := xml.Unmarshal(body, r)
		if e != nil {
			return e
		}
		return r
	}
	return nil
}

func (o *oss) do(method string, obj string, data io.Reader) (*http.Response, error) {
	if o.cli == nil {
		return nil, errors.New("no http client provided")
	}
	obj = strings.TrimLeft(obj, "/")
	url := "https://" + o.bucket + "." + o.endpoint + "/" + obj
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Set("Authorization", o.sign(req, obj))
	return o.cli.Do(req)
}

func (o *oss) sign(req *http.Request, obj string) string {
	md5 := req.Header.Get("Content-Md5")
	ctype := req.Header.Get("Content-Type")
	s := req.Method + "\n" + md5 + "\n" + ctype + "\n" + req.Header.Get("Date") + "\n" + o.canonicalize(req.Header) + "/" + o.bucket + "/"
	if obj != "" {
		s += obj
	}
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(o.as))
	_, _ = io.WriteString(h, s)
	return "OSS " + o.ak + ":" + base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (o *oss) canonicalize(header http.Header) string {
	mmap := make(map[string]string)
	kslice := make([]string, 0, len(header))
	for k, v := range header {
		lk := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(lk, "x-oss-") {
			kslice = append(kslice, lk)
			if len(v) == 0 {
				mmap[lk] = ""
			} else {
				mmap[lk] = strings.TrimSpace(v[0])
			}
		}
	}
	if len(kslice) == 0 {
		return ""
	}
	sort.Strings(kslice)
	result := ""
	for _, s := range kslice {
		result += s + ":" + mmap[s] + "\n"
	}
	return result
}
