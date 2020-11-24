package chlorine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.badanamu.com.cn/calmisland/common-cn/helper"

	"gitlab.badanamu.com.cn/calmisland/common-log/log"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) *Client {
	c := &Client{
		endpoint: endpoint,
	}
	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	return c
}

func (c *Client) Run(ctx context.Context, req *Request, resp *Response) (int, error) {
	reqBody := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     req.q,
		Variables: req.vars,
	}
	reqBuffer, err := json.Marshal(&reqBody)
	if err != nil {
		log.Warn(ctx, "Run: Marshal failed", log.Err(err), log.Any("reqBody", reqBody))
		return 0, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewBuffer(reqBuffer))
	if err != nil {
		log.Warn(ctx, "Run: New httpRequest failed", log.Err(err), log.Any("reqBody", reqBody))
		return 0, err
	}
	if bada, ok := helper.GetBadaCtx(ctx); ok {
		bada.SetHeader(request.Header)
	}
	request.Header = req.Header
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application; charset=utf-8")
	res, err := c.httpClient.Do(request)
	if err != nil {
		log.Error(ctx, "Run: do http failed", log.Err(err), log.String("endpoint", c.endpoint), log.Any("reqBody", reqBody))
		return 0, err
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(ctx, "Run: read response failed",
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return 0, err
	}
	err = json.Unmarshal(response, resp)
	if err != nil {
		log.Error(ctx, "Run: unmarshal response failed",
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return 0, err
	}
	return res.StatusCode, nil
}

type Request struct {
	q      string
	vars   map[string]interface{}
	Header http.Header
}

type OptFunc func(*Request)

func ReqToken(token string) OptFunc {
	return func(req *Request) {
		req.Header.Add(cookieKey, fmt.Sprintf("access=%s", token))
	}
}

func NewRequest(q string, opt ...OptFunc) *Request {
	req := &Request{
		q:      q,
		Header: make(map[string][]string),
	}

	for i := range opt {
		opt[i](req)
	}
	return req
}

func (req *Request) Var(key string, value interface{}) {
	if req.vars == nil {
		req.vars = make(map[string]interface{})
	}
	req.vars[key] = value
}

func (req *Request) SetHeader(key string, value string) {
	req.Header[key] = []string{value}
}

func (req *Request) SetHeaders(key string, values []string) {
	req.Header[key] = values
}

const cookieKey = "Cookie"

type ClError struct {
	Message   string `json:"message"`
	Locations []struct {
		Line   int `json:"line"`
		Column int `json:"column"`
	} `json:"locations"`
	Extensions struct {
		Code      string `json:"code"`
		Exception struct {
			Stacktrace []string `json:"stacktrace"`
		} `json:"exception"`
	}
}

type ClErrors []*ClError

func (clErrs ClErrors) Error() string {
	if len(clErrs) > 0 {
		return clErrs[0].Message
	}
	return "Empty ClErrors"
}

type Response struct {
	Data   interface{} `json:"data,omitempty"`
	Errors ClErrors    `json:"errors,omitempty"`
}
