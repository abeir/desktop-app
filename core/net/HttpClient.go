package net

import (
	"bytes"
	"github.com/abeir/desktop-app/core"
	"github.com/abeir/desktop-app/core/log"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type HttpMethod string

func (h *HttpMethod) ToString() string{
	return string(*h)
}

const (
	HttpGet HttpMethod     = "GET"
	HttpHead HttpMethod    = "HEAD"
	HttpPost HttpMethod    = "POST"
	HttpPut HttpMethod     = "PUT"
	HttpPatch HttpMethod   = "PATCH"
	HttpDelete HttpMethod  = "DELETE"
	HttpConnect HttpMethod = "CONNECT"
	HttpOptions HttpMethod = "OPTIONS"
	HttpTrace HttpMethod   = "TRACE"
)

type ContentType string

func (c *ContentType) ToString() string{
	return string(*c)
}
const (
	FormUrlencoded ContentType = "application/x-www-form-urlencoded"
	FormData ContentType       = "multipart/form-data"
	Json ContentType           = "application/json"
)


// NewHttpClient 创建HttpClient实例
// HttpClient中的Set、Add方法支持链式调用方式，最后通过Reqeust方法发送请求
// 默认情况下，会将Content-Type设置为application/x-www-form-urlencoded
func NewHttpClient() *HttpClient{
	client := fasthttp.Client{
		//每个主机的最大连接数
		MaxConnsPerHost: 50,
		//空闲的连接关闭时间
		MaxIdleConnDuration: 10 * time.Second,
		//活动的连接关闭时间
		MaxConnDuration: 60 * time.Second,
		//幂等调用的最大尝试次数
		MaxIdemponentCallAttempts: 3,
		//响应数据的读取超时时间
		ReadTimeout: 8 * time.Second,
		//请求数据的写入超时时间
		WriteTimeout: 8 * time.Second,
	}
	httpClient := &HttpClient{
		client: &client,
		method: HttpGet,
		req: fasthttp.AcquireRequest(),
		respHeaders: make(map[string][]byte),
	}
	httpClient.SetContentType(FormUrlencoded)
	return httpClient
}

// HttpClient http客户端，不要手动构建该实例，应该调用NewHttpClient函数创建实例
type HttpClient struct {
	client *fasthttp.Client
	method HttpMethod
	req *fasthttp.Request
	respHeaders map[string][]byte
	err error
}

// AddHeader 添加请求头
func (h *HttpClient) AddHeader(name, value string) *HttpClient{
	header := h.req.Header
	header.Add(name, value)
	return h
}

// AddHeaders 批量添加请求头
func (h *HttpClient) AddHeaders(headers map[string]string) *HttpClient{
	header := h.req.Header
	for k,v := range headers {
		header.Add(k, v)
	}
	return h
}

// SetContentType 设置Content-Type
func (h *HttpClient) SetContentType(contentType ContentType) *HttpClient{
	h.req.Header.SetContentType(contentType.ToString())
	return h
}

// SetUserAgent 设置User-Agent
func (h *HttpClient) SetUserAgent(userAgent string) *HttpClient{
	h.req.Header.SetUserAgent(userAgent)
	return h
}

// SetMethod 设置请求方法，GET、POST、DELETE等
func (h *HttpClient) SetMethod(method HttpMethod) *HttpClient{
	h.method = method
	h.req.Header.SetMethod(h.method.ToString())
	return h
}

// SetBody 设置请求体内容
func (h *HttpClient) SetBody(body []byte) *HttpClient{
	h.req.SetBody(body)
	return h
}

func (h *HttpClient) SetBodyStream(bodyStream io.Reader) *HttpClient{
	h.req.SetBodyStream(bodyStream, -1)
	return h
}

func (h *HttpClient) SetBodyMap(body map[string][]string) *HttpClient{
	if body==nil || len(body)==0 {
		return h
	}
	var bodyBuff bytes.Buffer
	for name,values := range body {
		if values==nil {
			bodyBuff.WriteString(name)
			bodyBuff.WriteString("=")
			bodyBuff.WriteString("&")
			continue
		}
		for _,val := range values {
			bodyBuff.WriteString(name)
			bodyBuff.WriteString("=")
			bodyBuff.WriteString(val)
			bodyBuff.WriteString("&")
		}
	}
	if bodyBuff.Len() == 0 {
		return h
	}
	bodyBytes := bodyBuff.Bytes()
	if bodyBytes[len(bodyBytes)-1] == '&' {
		bodyBytes = bodyBytes[:len(bodyBytes)-1]
	}
	h.SetBody(bodyBytes)
	return h
}


// MultipartForm 用于发送multipart/form-data类型的数据，会将Content-Type设置为multipart/form-data
//    params: 请求参数
//    files: 上传文件
func (h *HttpClient) MultipartForm(params map[string][]string, files map[string]string) *HttpClient{
	var bodyReader bytes.Buffer
	bodyWriter := multipart.NewWriter(&bodyReader)

	if params!=nil {
		for name,values := range params {
			if values==nil || len(values)==0 {
				if err := bodyWriter.WriteField(name, ""); err!=nil {
					h.err = err
					log.Error(err)
					return h
				}
				continue
			}
			for _, val := range values {
				if err := bodyWriter.WriteField(name, val); err!=nil {
					h.err = err
					log.Error(err)
					return h
				}
			}
		}
	}
	if files!=nil {
		for name,file := range files {
			filename := filepath.Base(file)
			fileWriter, err := bodyWriter.CreateFormFile(name, filename)
			if err!=nil {
				h.err = err
				log.Error(err)
				return h
			}
			if err = h.copy(file, fileWriter); err!=nil {
				h.err = err
				log.Error(err)
				return h
			}
		}
	}
	if err := bodyWriter.Close(); err!=nil {
		h.err = err
		log.Error(err)
		return h
	}
	h.SetContentType(ContentType(bodyWriter.FormDataContentType()))
	h.SetBodyStream(&bodyReader)
	return h
}

func (h *HttpClient) copy(srcFile string, dst io.Writer) error{
	f, _ := os.Open(srcFile)
	defer core.CloseQuietly(f)
	if _, err := io.Copy(dst, f); err!=nil {
		return err
	}
	return nil
}

func (h *HttpClient) Request(url string) (body []byte, err error){
	if h.err!=nil {
		return nil, h.err
	}
	h.req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(h.req)
		h.req = nil
		fasthttp.ReleaseResponse(resp)
	}()
	if err = h.client.Do(h.req, resp); err!=nil {
		return nil, err
	}
	resp.Header.VisitAll(func(key, value []byte){
		h.respHeaders[string(key)] = value
	})
	body = resp.Body()
	return body, nil
}

func (h *HttpClient) RequestStream(url string) (body io.Writer, err error){
	if h.err!=nil {
		return nil, h.err
	}
	h.req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(h.req)
		h.req = nil
		fasthttp.ReleaseResponse(resp)
	}()
	if err = h.client.Do(h.req, resp); err!=nil {
		return nil, err
	}
	resp.Header.VisitAll(func(key, value []byte){
		h.respHeaders[string(key)] = value
	})
	body = resp.BodyWriter()
	return body, nil
}

func (h *HttpClient) ResponseHeaders() map[string][]byte{
	return h.respHeaders
}


func (h *HttpClient) FastGet(url string) (body []byte, err error){
	var result []byte
	_, body, err = h.client.Get(result, url)
	return body, err
}

func (h *HttpClient) FastPost(url string, data map[string]string) (body []byte, err error){
	var result []byte

	args := fasthttp.AcquireArgs()
	for k,v := range data {
		args.Add(k, v)
	}
	defer fasthttp.ReleaseArgs(args)
	_, body, err = h.client.Post(result, url, args)
	return body, err
}
