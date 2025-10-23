package chlorine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"
	"time"

	"github.com/Klasmart-Engineering/common-log/log"
	"github.com/Klasmart-Engineering/tracecontext"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Client struct {
	endpoint    string
	httpClient  *http.Client
	httpTimeout time.Duration
}

const defaultHttpTimeout = time.Minute

type OptionChlorine func(c *Client)

func WithTimeout(duration time.Duration) OptionChlorine {
	return func(c *Client) {
		c.httpTimeout = duration
	}
}

func DisableNewRelicDistributedTracing(httpClient *http.Client) OptionChlorine {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

type debugTransport struct{}

func (d debugTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx := req.Context()
	txnExist := newrelic.FromContext(ctx) != nil
	log.Debug(ctx, "chlorine costume round trip",
		log.Any("headers", req.Header),
		log.Bool("txn exist", txnExist))
	return http.DefaultTransport.RoundTrip(req)
}

func NewClient(endpoint string, options ...OptionChlorine) *Client {
	c := &Client{
		endpoint: endpoint,
		httpClient:  &http.Client{Transport: newrelic.NewRoundTripper(debugTransport{})},
		httpTimeout: defaultHttpTimeout,
	}
	for i := range options {
		options[i](c)
	}
	return c
}

func (c *Client) Run(ctx context.Context, req *Request, resp *Response) (int, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.httpTimeout)
	defer cancel()

	reqBody := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     req.q,
		Variables: req.vars,
	}
	reqBuffer, err := json.Marshal(&reqBody)
	if err != nil {
		log.Warn(ctxWithTimeout, "Run: Marshal failed", log.Err(err), log.Any("reqBody", reqBody))
		return 0, err
	}
	request, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, c.endpoint, bytes.NewBuffer(reqBuffer))
	if err != nil {
		log.Warn(ctxWithTimeout, "Run: New httpRequest failed", log.Err(err), log.Any("reqBody", reqBody))
		return 0, err
	}
	if bada, ok := tracecontext.GetTraceContext(ctxWithTimeout); ok {
		bada.SetHeader(request.Header)
	}
	request.Header = req.Header
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept", "application; charset=utf-8")
	var result *http.Response
	var resultErr error

	start := time.Now()
	result, resultErr = c.httpClient.Do(request)

	duration := time.Since(start)
	if resultErr != nil {
		log.Error(ctxWithTimeout, "Run: do http failed",
			log.Duration("duration", duration),
			log.Err(resultErr),
			log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody))
		return 0, resultErr
	}

	defer result.Body.Close()
	response, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error(ctxWithTimeout, "Run: read response failed",
			log.Duration("duration", duration),
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return result.StatusCode, err
	}

	err = json.Unmarshal(response, resp)
	if err != nil {
		log.Error(ctxWithTimeout, "Run: unmarshal response failed",
			log.Duration("duration", duration),
			log.Err(err), log.String("endpoint", c.endpoint),
			log.Any("reqBody", reqBody), log.String("response", string(response)))
		return result.StatusCode, err
	}
	log.Debug(ctxWithTimeout, "Run: Success",
		log.Duration("duration", duration),
		log.Any("reqBody", reqBody),
		log.String("response", string(response)))
	return result.StatusCode, nil
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
