package repo

import (
	"context"
	cos "github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"os"
)

type Cos struct {
	c       *cos.Client
	limiter *rate.Limiter
}

type CosOption struct {
	SecretIdValue  string
	SecretKeyValue string
	BucketURL      string
	ServiceURL     string
}

func NewCos(conf CosOption) (*Cos, error) {
	bu, err := url.Parse(conf.BucketURL)
	if err != nil {
		return nil, err
	}

	rates := rate.Limit(1.0 / 3)
	limiter := rate.NewLimiter(rates, 1)

	su, err := url.Parse(conf.ServiceURL)
	if err != nil {
		return nil, err
	}

	b := &cos.BaseURL{BucketURL: bu, ServiceURL: su}

	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretIdValue,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: conf.SecretKeyValue, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})

	return &Cos{
		c:       c,
		limiter: limiter,
	}, nil
}

func (cos *Cos) Process(ctx context.Context, cosFileName, localFilePath string) error {
	f, err := os.Open(localFilePath)
	if err != nil {
		return err
	}

	_, err = cos.c.Object.Put(ctx, cosFileName, f, nil)

	return err
}
