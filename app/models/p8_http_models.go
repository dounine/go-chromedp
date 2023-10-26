package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	P8Http struct {
	}
	P8RequestHeader struct {
		headers map[string]string
	}
	P8RequestParam struct {
		params map[string]string
	}
	P8Request struct {
		request *http.Request
		method  string
		url     string
		Header  P8RequestHeader
		Param   P8RequestParam
		ctx     *context.Context
		Body    any
	}
)

func NewP8Http() *P8Http {
	return &P8Http{}
}
func (c *P8RequestHeader) Add(key string, value string) {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[key] = value
}

func (c *P8RequestParam) Add(key string, value string) {
	if c.params == nil {
		c.params = make(map[string]string)
	}
	c.params[key] = value
}

func (c *P8Http) Get(ctx *context.Context, url string) *P8Request {
	request := &P8Request{
		method: http.MethodGet,
		url:    url,
		ctx:    ctx,
	}
	return request
}
func (c *P8Http) Post(ctx *context.Context, url string) *P8Request {
	request := &P8Request{
		method: http.MethodPost,
		url:    url,
		ctx:    ctx,
	}
	return request
}
func (c *P8Http) Patch(ctx *context.Context, url string) *P8Request {
	request := &P8Request{
		method: http.MethodPatch,
		url:    url,
		ctx:    ctx,
	}
	return request
}
func (c *P8Http) Put(ctx *context.Context, url string) *P8Request {
	request := &P8Request{
		method: http.MethodPut,
		url:    url,
		ctx:    ctx,
	}
	return request
}
func (c *P8Http) Delete(ctx *context.Context, url string) *P8Request {
	request := &P8Request{
		method: http.MethodDelete,
		url:    url,
		ctx:    ctx,
	}
	return request
}

func (c *P8Request) get() (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(*c.ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, nil, err
	}
	return c.response(req)
}

func (c *P8Request) post() (*http.Response, []byte, error) {
	var reader = io.Reader(nil)
	if c.Body != nil {
		data, err := json.Marshal(c.Body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewBuffer(data)
	}
	req, err := http.NewRequestWithContext(*c.ctx, http.MethodPost, c.url, reader)
	if err != nil {
		return nil, nil, err
	}
	return c.response(req)
}

func (c *P8Request) requestReader() io.Reader {
	var reader = io.Reader(nil)
	if c.Body != nil {
		data, err := json.Marshal(c.Body)
		if err != nil {
			return nil
		}
		reader = bytes.NewReader(data)
	}
	return reader
}

func (c *P8Request) response(req *http.Request) (*http.Response, []byte, error) {
	client := http.DefaultClient
	if c.Header.headers != nil {
		for key, value := range c.Header.headers {
			req.Header.Add(key, value)
		}
	}
	if c.Param.params != nil {
		q := req.URL.Query()
		for key, value := range c.Param.params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			fmt.Println(closeError)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, body, nil
}

func (c *P8Request) patch() (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(*c.ctx, http.MethodPatch, c.url, c.requestReader())
	if err != nil {
		return nil, nil, err
	}
	return c.response(req)
}

func (c *P8Request) put() (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(*c.ctx, http.MethodPut, c.url, c.requestReader())
	if err != nil {
		return nil, nil, err
	}
	return c.response(req)
}

func (c *P8Request) delete() (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(*c.ctx, http.MethodDelete, c.url, nil)
	if err != nil {
		return nil, nil, err
	}
	return c.response(req)
}

func (c *P8Request) End() (*http.Response, []byte, error) {
	switch c.method {
	case http.MethodGet:
		return c.get()
	case http.MethodPost:
		return c.post()
	case http.MethodPatch:
		return c.patch()
	case http.MethodPut:
		return c.put()
	case http.MethodDelete:
		return c.delete()
	}
	return nil, nil, nil
}
