package aliyun

import (
	"errors"
	"net/http"
	"strings"
)

type oss struct {
	ak       string
	as       string
	bucket   string
	endpoint string
	client   *http.Client
}

func NewOssClient(ak, as, bucket, endpoint string, cli *http.Client) *oss {
	endpoint = strings.TrimRight(endpoint, "/")
	return &oss{
		ak:       ak,
		as:       as,
		bucket:   bucket,
		endpoint: endpoint,
		client:   cli,
	}
}

func (o *oss) PutObject(bs []byte, obj string) error {
	if o.client == nil {
		return errors.New("no http client provided")
	}
	//remote := "https://" + o.bucket + "." + o.endpoint
	return nil
}

func (o *oss) DeleteObject(obj string) {

}
