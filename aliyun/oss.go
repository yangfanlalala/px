package aliyun

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
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
	if resp.StatusCode != http.StatusOK {
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			return e
		}
		r := &errResponse{}
		e = json.Unmarshal(body, r)
		if e != nil {
			return err
		}
		return r
	}
	return nil
}

func (o *oss) GetObjectURL() string {
	return ""
}

func (o *oss) DeleteObject() error {
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
	return o.cli.Do(req)
}

func (o *oss) sign(req *http.Request, obj string) string {
	date := time.Now().Format(http.TimeFormat)
	req.Header.Set("Date", date)
	s := req.Method + "\n\n\n" + date + o.canonicalize(req.Header) + "/" + o.bucket
	if obj != "" {
		s += "/" + obj
	}
	return "OSS" + o.ak
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
