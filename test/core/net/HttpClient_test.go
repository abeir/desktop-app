package net

import (
	"bytes"
	"github.com/abeir/desktop-app/core/net"
	"testing"
)

func TestRequestGet(t *testing.T) {
	client := net.NewHttpClient()
	client.SetMethod(net.HttpGet)
	body, e := client.Request("http://127.0.0.1:8000/get?test=123")
	if e!=nil {
		t.Error(e)
	}
	if body==nil || len(body)==0 {
		t.Error("响应内容为空")
	}
}

func TestRequestPost1(t *testing.T){
	client := net.NewHttpClient()
	client.SetMethod(net.HttpPost)
	client.AddHeader("user", "1")
	client.SetBody([]byte("body content!!!"))
	body, e := client.Request("http://127.0.0.1:8000/post")
	if e!=nil {
		t.Error(e)
	}
	if body==nil || len(body)==0 {
		t.Error("响应内容为空")
	}
}

func TestRequestPost2(t *testing.T){
	client := net.NewHttpClient()
	client.SetMethod(net.HttpPost)
	client.AddHeader("user", "2")

	var b bytes.Buffer
	b.WriteString("post2")
	client.SetBodyStream(&b)

	body, e := client.Request("http://127.0.0.1:8000/post")
	if e!=nil {
		t.Error(e)
	}
	if body==nil || len(body)==0 {
		t.Error("响应内容为空")
	}
}

func TestRequestPost3(t *testing.T){
	client := net.NewHttpClient()
	client.SetMethod(net.HttpPost)
	client.AddHeader("user", "3")

	m := make(map[string][]string)
	m["name"] = []string{"abeir"}
	m["age"] = []string{"23"}
	client.SetBodyMap(m)

	body, e := client.Request("http://127.0.0.1:8000/post")
	if e!=nil {
		t.Error(e)
	}
	if body==nil || len(body)==0 {
		t.Error("响应内容为空")
	}
}