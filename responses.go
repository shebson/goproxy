package goproxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Will generate a valid http response to the given request the response will have
// the given contentType, and http status.
// Typical usage, refuse to process requests to local addresses:
//
//	proxy.OnRequest(IsLocalHost()).DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request,*http.Response) {
//		return nil,NewResponse(r,goproxy.ContentTypeHtml,http.StatusUnauthorized,
//			`<!doctype html><html><head><title>Can't use proxy for local addresses</title></head><body/></html>`)
//	})
func NewResponse(r *http.Request, contentType string, status int, body string) *http.Response {
	resp := &http.Response{}
	resp.Request = r
	resp.TransferEncoding = r.TransferEncoding
	resp.Header = make(http.Header)
	resp.Header.Add("Content-Type", contentType)
	resp.StatusCode = status
	buf := bytes.NewBufferString(body)
	resp.ContentLength = int64(buf.Len())
	resp.Body = ioutil.NopCloser(buf)
	return resp
}

func NewForbiddenResponse(req *http.Request, realm string) *http.Response {
	var unauthorizedMsg = []byte("407 Proxy Authentication Required")
	var response = &http.Response{
		StatusCode:    407,
		ProtoMajor:    1,
		ProtoMinor:    1,
		Request:       req,
		Header:        http.Header{
			"Proxy-Authenticate": []string{"Basic realm=" + realm},
			"Connection": []string{"close"},
			},
		Body:          ioutil.NopCloser(bytes.NewBuffer(unauthorizedMsg)),
		ContentLength: int64(len(unauthorizedMsg)),
	}
	return response
}

const (
	ContentTypeText = "text/plain"
	ContentTypeHtml = "text/html"
)

// Alias for NewResponse(r,ContentTypeText,http.StatusAccepted,text)
func TextResponse(r *http.Request, text string) *http.Response {
	return NewResponse(r, ContentTypeText, http.StatusAccepted, text)
}
