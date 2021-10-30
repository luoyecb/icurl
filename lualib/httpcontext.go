package lualib

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"
)

type HttpContext struct {
	Scheme string
	Host   string
	Port   int
	Path   string
	Method string
	Url    string // if set, other fields are ignored.
	Data   string // if Method=GET, Data are ignored, otherwise use Data instead of Query
	Query  map[string]string
	Header map[string]string
}

func NewHttpContext() *HttpContext {
	ctx := new(HttpContext)
	return ctx
}

func (ctx *HttpContext) Send() (string, error) {
	url := ctx.BuildUrl()
	if url == "" {
		return "", errors.New("http context info invalid")
	}

	request := gorequest.New()

	method := strings.ToUpper(ctx.Method)
	if method == "GET" {
		request.Get(url)
	} else if method == "POST" {
		request.Post(url)
	} else if method == "PUT" {
		request.Put(url)
	} else if method == "DELETE" {
		request.Delete(url)
	} else {
		return "", errors.New("only supported GET|POST|PUT|DELETE method")
	}

	if method != "GET" {
		if ctx.Data != "" {
			request.Send(ctx.Data)
		} else {
			request.SendMap(ctx.Query)
		}
	}

	// Note：设置header的代码必须放在 request.Get()|Post()|Delete()|Put() 之后
	if len(ctx.Header) > 0 {
		for hk, hv := range ctx.Header {
			request.Set(hk, hv)
		}
	}

	// Before request
	fmt.Printf("=== Send request to (%s)%s\n", method, url)

	resp, bodyStr, errs := request.End()
	if len(errs) > 0 {
		return "", errs[0]
	}

	// After request
	fmt.Printf("=== Status code: %d\n", resp.StatusCode)
	fmt.Printf("=== Response header\n")
	for k, v := range resp.Header {
		fmt.Printf("%s = %s\n", k, v[0])
	}
	fmt.Println()

	return bodyStr, nil
}

func (ctx *HttpContext) BuildUrl() string {
	// first use Url
	if ctx.Url != "" {
		return ctx.Url
	}

	var buf bytes.Buffer

	// Scheme
	if ctx.Scheme == "https" || ctx.Scheme == "http" {
		buf.WriteString(ctx.Scheme)
	} else {
		return ""
	}

	// Host
	if ctx.Host != "" {
		buf.WriteString("://")
		buf.WriteString(ctx.Host)
	} else {
		return ""
	}

	// Port
	if ctx.Port > 0 && ctx.Port != 80 && ctx.Port != 443 {
		buf.WriteString(fmt.Sprintf(":%d", ctx.Port))
	}

	// Path
	path := ctx.Path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	}
	buf.WriteString(path)

	// Query string
	if CompareStringIgnoreCase("get", ctx.Method) && len(ctx.Query) > 0 {
		query := make(map[string][]string)
		for k, v := range ctx.Query {
			query[k] = append(query[k], v)
		}
		queryStr := url.Values(query).Encode()
		buf.WriteString("?")
		buf.WriteString(queryStr)
	}

	return buf.String()
}
