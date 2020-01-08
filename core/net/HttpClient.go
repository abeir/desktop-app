package net

import (
	"bufio"
	"bytes"
	"github.com/abeir/desktop-app/core"
	"github.com/abeir/desktop-app/core/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HttpMethod http请求类型，post、get等等
type HttpMethod string

func (h HttpMethod) ToString() string{
	return string(h)
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

// ContentType 请求头Content-Type的值
type ContentType string

func (c ContentType) ToString() string{
	return string(c)
}
const (
	FormUrlencoded ContentType = "application/x-www-form-urlencoded"
	FormData ContentType       = "multipart/form-data"
	Json ContentType           = "application/json"
)

var client *http.Client

// NewHttpClient 创建HttpClient实例
// HttpClient中的Set、Add方法支持链式调用方式，最后通过Reqeust方法发送请求
// 默认情况下，会将Content-Type设置为application/x-www-form-urlencoded
func NewHttpClient() *HttpClient{
	if client==nil {
		var lock = sync.Mutex{}
		lock.Lock()
		defer lock.Unlock()
		if client==nil {
			var transport = &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			}
			client = &http.Client{
				Transport:     transport,
				Timeout:       time.Second * 15,
			}
		}
	}

	httpClient := &HttpClient{
		client: client,
		method: HttpGet,
		headers: make(map[string]string),
		rspHeaders: make(map[string][]string),
	}
	httpClient.SetContentType(FormUrlencoded)
	return httpClient
}

// HttpClient http客户端，不要手动构建该实例，应该调用NewHttpClient函数创建实例
type HttpClient struct {
	client *http.Client
	method HttpMethod
	headers map[string]string
	cookies []http.Cookie
	body io.Reader
	rspHeaders map[string][]string
	err error
}

// AddHeader 添加请求头
func (h *HttpClient) AddHeader(name, value string) *HttpClient{
	h.headers[name] = value
	return h
}

// AddHeaders 批量添加请求头
func (h *HttpClient) AddHeaders(headers map[string]string) *HttpClient{
	for k,v := range headers {
		h.headers[k] = v
	}
	return h
}

func (h *HttpClient) AddCookie(cookie *http.Cookie) *HttpClient{
	h.cookies = append(h.cookies, *cookie)
	return h
}

// SetContentType 设置Content-Type
func (h *HttpClient) SetContentType(contentType ContentType) *HttpClient{
	h.headers["Content-Type"] = contentType.ToString()
	return h
}

// SetUserAgent 设置User-Agent
func (h *HttpClient) SetUserAgent(userAgent string) *HttpClient{
	h.headers["User-Agent"] = userAgent
	return h
}

// SetMethod 设置请求方法，GET、POST、DELETE等
func (h *HttpClient) SetMethod(method HttpMethod) *HttpClient{
	h.method = method
	return h
}

// SetBody 设置请求体内容
func (h *HttpClient) SetBody(body []byte) *HttpClient{
	h.body = bytes.NewBuffer(body)
	return h
}

// SetBodyStream 设置请求体内容
func (h *HttpClient) SetBodyStream(body io.Reader) *HttpClient{
	h.body = body
	return h
}

// SetBodyMap 设置请求体内容，通常post请求参数可以放置在此处
func (h *HttpClient) SetBodyMap(body map[string][]string) *HttpClient{
	bodyBytes := h.bodyMap2bytes(body)
	if len(bodyBytes)==0 {
		return h
	}
	h.SetBody(bodyBytes)
	return h
}

func (h *HttpClient) bodyMap2bytes(body map[string][]string) []byte{
	if body==nil || len(body)==0 {
		return []byte{}
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
		return []byte{}
	}
	bodyBytes := bodyBuff.Bytes()
	if bodyBytes[len(bodyBytes)-1] == '&' {
		bodyBytes = bodyBytes[:len(bodyBytes)-1]
	}
	return bodyBytes
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

func (h *HttpClient) doRequest(url string) (rsp *http.Response, err error){
	if h.err!=nil {
		return nil, h.err
	}
	req, err := http.NewRequest(h.method.ToString(), url, h.body)
	if err!=nil {
		return nil, err
	}
	for k,v := range h.headers {
		req.Header.Add(k, v)
	}
	for _, cookie := range h.cookies {
		req.AddCookie(&cookie)
	}
	return h.client.Do(req)
}

func (h *HttpClient) extractRspHeaders(rsp *http.Response){
	for k,v := range rsp.Header {
		h.rspHeaders[k] = v
	}
}

// Request 发送请求，调用成功后可使用 ResponseHeaders方法获取响应头
//    url: 请求地址
// return
//    body: 响应内容
//    err: 请求过程中出现的错误
func (h *HttpClient) Request(url string) (body []byte, err error){
	rsp, err := h.doRequest(url)
	if err!=nil {
		return nil, err
	}
	h.extractRspHeaders(rsp)
	defer core.CloseQuietly(rsp.Body)
	body, err = ioutil.ReadAll(rsp.Body)
	return body, err
}

// RequestStream 发送请求，调用成功后可使用 ResponseHeaders方法获取响应头
//    url: 请求地址
// return
//    body: 响应内容
//    err: 请求过程中出现的错误
func (h *HttpClient) RequestStream(url string) (body io.Writer, err error){
	rsp, err := h.doRequest(url)
	if err!=nil {
		return nil, err
	}
	h.extractRspHeaders(rsp)
	defer core.CloseQuietly(rsp.Body)

	var bodyWriter bufio.Writer
	_, err = io.Copy(&bodyWriter, rsp.Body)
	return body, err
}

// ResponseHeaders 请求完成后，使用此方法获取响应头
func (h *HttpClient) ResponseHeaders() map[string][]string{
	return h.rspHeaders
}

// FastGet 发送简单的GET请求，注意，调用该方法发送请求后 ResponseHeaders方法不会获取响应头
//    url: 请求地址
// return
//    body: 响应内容
//    err: 请求过程中出现的错误
func (h *HttpClient) FastGet(url string) (body []byte, err error){
	rsp, err := h.client.Get(url)
	if err!=nil {
		return nil, err
	}
	defer core.CloseQuietly(rsp.Body)
	body, err = ioutil.ReadAll(rsp.Body)
	return body, err
}

// FastPost 发送简单的POST请求，注意，调用该方法发送请求后 ResponseHeaders方法不会获取响应头
//    url: 请求地址
//    data: post请求参数
// return
//    body: 响应内容
//    err: 请求过程中出现的错误
func (h *HttpClient) FastPost(url string, data map[string][]string) (body []byte, err error){
	requestBody := h.bodyMap2bytes(data)
	rsp, err := h.client.Post(url, FormUrlencoded.ToString(), bytes.NewBuffer(requestBody))
	if err!=nil {
		return nil, err
	}
	defer core.CloseQuietly(rsp.Body)
	body, err = ioutil.ReadAll(rsp.Body)
	return body, err
}
