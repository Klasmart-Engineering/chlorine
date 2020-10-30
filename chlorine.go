package chlorine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func (c *Client) Run(ctx context.Context, req *Request, resp interface{}) error {
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
		return err
	}
	request, err := http.NewRequest(http.MethodPost, c.endpoint, bytes.NewBuffer(reqBuffer))
	if err != nil {
		log.Warn(ctx, "Run: New httpRequest failed", log.Err(err), log.Any("reqBody", reqBody))
		return err
	}
	if bada, ok := helper.GetBadaCtx(ctx); ok {
		bada.SetHeader(request.Header)
	}
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application; charset=utf-8")
	for key, values := range req.Header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}
	res, err := c.httpClient.Do(request)
	if err != nil {
		log.Error(ctx, "Run: do http failed", log.Err(err), log.String("endpoint", c.endpoint), log.Any("reqBody", reqBody))
		return err
	}
	defer res.Body.Close()
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(ctx, "Run: read response failed",
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return err
	}
	if res.StatusCode != http.StatusOK {
		log.Warn(ctx, "Run: read response failed",
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return errors.New(res.Status + string(response))
	}
	err = json.Unmarshal(response, resp)
	if err != nil {
		log.Error(ctx, "Run: unmarshal response failed",
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return err
	}
	return nil
}

type Request struct {
	q      string
	vars   map[string]interface{}
	Header http.Header
}

func NewRequest(q string) *Request {
	req := &Request{
		q:      q,
		Header: make(map[string][]string),
	}
	return req
}

func (req *Request) Var(key string, value interface{}) {
	if req.vars == nil {
		req.vars = make(map[string]interface{})
	}
	req.vars[key] = value
}

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
