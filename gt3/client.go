package gt3

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nilorg/geetest/pkg/util"
)

// Client gt3客户端
type Client struct {
	geetestID  string // 公钥
	geetestKey string // 私钥
	opts       *ClientOptions
	httpClient *http.Client
}

// NewClient 创建客户端
func NewClient(geetestID string, geetestKey string, opts ...Option) *Client {
	o := newOptions(opts...)
	return &Client{
		geetestID:  geetestID,
		geetestKey: geetestKey,
		opts:       &o,
		httpClient: &http.Client{Timeout: time.Duration(o.HTTPTimeout) * time.Second},
	}
}

func (g *Client) buildRequestComm(userID string) *RequestComm {
	return &RequestComm{
		UserID:     userID,
		ClientType: g.opts.ClientType,
		IPAddress:  g.opts.IPAddress,
		JSONFormat: g.opts.JSONFormat,
		Sdk:        g.opts.Version,
	}
}

// httpGet 发送GET请求，获取服务器返回结果
func (g *Client) httpGet(uri string, params map[string]string) (body []byte, err error) {
	q := url.Values{}
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?%s", g.opts.APIURL, uri, q.Encode()), nil)
	if err != nil {
		return
	}
	var res *http.Response
	res, err = g.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = fmt.Errorf("response status code: %d", res.StatusCode)
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	return
}

// httpPost 发送POST请求，获取服务器返回结果
func (g *Client) httpPost(uri string, params map[string]string) (body []byte, err error) {
	q := url.Values{}
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", g.opts.APIURL, uri), strings.NewReader(q.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var res *http.Response
	res, err = g.httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = fmt.Errorf("response status code: %d", res.StatusCode)
		return
	}
	body, err = ioutil.ReadAll(res.Body)
	return
}

// Register register
func (g *Client) Register(digestmod string, userID ...string) (res *RegisterResponse, err error) {
	var reqComm *RequestComm
	if len(userID) > 0 {
		reqComm = g.buildRequestComm(userID[0])
	} else {
		reqComm = g.buildRequestComm("")
	}
	req := &RegisterRequest{
		RequestComm: reqComm,
		Gt:          g.geetestID,
		Digestmod:   digestmod,
	}
	params := util.StructToMap(req)
	var body []byte
	body, err = g.httpGet(g.opts.RegisterURL, params)
	if err != nil {
		return
	}
	res = new(RegisterResponse)
	err = json.Unmarshal(body, res)
	if err != nil {
		res = nil
	}
	return
}

// Validate validate
func (g *Client) Validate(challenge string, userID ...string) (res *ValidateResponse, err error) {
	var reqComm *RequestComm
	if len(userID) > 0 {
		reqComm = g.buildRequestComm(userID[0])
	} else {
		reqComm = g.buildRequestComm("")
	}
	req := &ValidateRequest{
		RequestComm: reqComm,
		CaptchaID:   g.geetestID,
		Seccode:     fmt.Sprintf("%s|jordan", g.geetestKey),
		Challenge:   challenge,
	}
	params := util.StructToMap(req)
	var body []byte
	body, err = g.httpPost(g.opts.ValidteURL, params)
	if err != nil {
		return
	}
	res = new(ValidateResponse)
	err = json.Unmarshal(body, res)
	if err != nil {
		res = nil
	}
	return
}
