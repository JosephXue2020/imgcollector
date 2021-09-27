package downloader

import (
	"net/http"
	"time"
)

type Header struct {
	UserAgent string
	Cookie    string
	Referer   string
	Token     string
}

func setHeader(header *Header, req *http.Request) {
	if header.UserAgent != "" {
		req.Header.Set("User-Agent", header.UserAgent)
	}
	if header.Cookie != "" {
		req.Header.Set("Cookie", header.Cookie)
	}
	if header.Referer != "" {
		req.Header.Set("Referer", header.Referer)
	}
	if header.Token != "" {
		req.Header.Set("Token", header.Token)
	}

}

func Get(url string, header *Header) (*http.Response, error) {
	if header == nil {
		return http.Get(url)
	} else {
		client := &http.Client{}
		req, _ := http.NewRequest("Get", url, nil)
		setHeader(header, req)
		return client.Do(req)
	}
}

func SlowGet(url string, header *Header, sleepMilliSec int) (*http.Response, error) {
	resp, err := Get(url, header)

	// 设置默认睡眠时间500毫秒
	if sleepMilliSec == 0 {
		sleepMilliSec = 500
	}
	time.Sleep(time.Duration(sleepMilliSec * 1000000))

	return resp, err
}
