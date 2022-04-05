package lualib

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

type HttpContext struct {
	Url    string
	Method string
	Data   string // if data is not empty, use data
	Query  map[string]string
	Header map[string]string
}

func NewHttpContext() *HttpContext {
	return &HttpContext{}
}

func (ctx *HttpContext) buildUrl() string {
	if ctx.Url == "" {
		return ""
	}
	if strings.ToUpper(ctx.Method) == "GET" && len(ctx.Query) > 0 {
		val := url.Values{}
		for k, v := range ctx.Query {
			val.Add(k, v)
		}
		url := val.Encode()
		if strings.Contains(ctx.Url, "?") {
			return ctx.Url + "&" + url
		} else {
			return ctx.Url + "?" + url
		}
	}
	return ctx.Url
}

func (ctx *HttpContext) Send() (string, error) {
	url := ctx.buildUrl()
	if url == "" {
		return "", errors.New("http context info invalid")
	}

	request := gorequest.New().Timeout(3 * time.Second)

	method := strings.ToUpper(ctx.Method)
	switch method {
	case "GET":
		request.Get(url)
	case "POST":
		request.Post(url)
	case "PUT":
		request.Put(url)
	case "DELETE":
		request.Delete(url)
	default:
		return "", errors.New("only supported GET|POST|PUT|DELETE method")
	}

	if method != "GET" {
		if ctx.Data != "" {
			request.Send(ctx.Data)
		} else {
			request.SendMap(ctx.Query)
		}
	}

	if len(ctx.Header) > 0 {
		for hk, hv := range ctx.Header {
			request.Set(hk, hv)
		}
	}

	fmt.Printf("=== Send request to (%s)%s\n", method, url)
	resp, bodyStr, errs := request.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	fmt.Printf("=== Status code: %d\n", resp.StatusCode)
	fmt.Printf("=== Response header\n")
	for k, v := range resp.Header {
		fmt.Printf("%s = %s\n", k, v[0])
	}
	fmt.Println()

	return bodyStr, nil
}
