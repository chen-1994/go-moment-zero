package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
)

// Go 语言标准库中的net/http包十分的优秀，提供了非常完善的 HTTP 客户端与服务端的实现，仅通过几行代码就可以搭建一个非常简单的 HTTP 服务器。
// 几乎所有的 go 语言中的 web 框架，都是对已有的 http 包做的封装与修改，因此，十分建议学习其他框架前先行掌握 http 包
// Get
func Test1() {
	resp, err := http.Get("http://www.baidu.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	content, _ := io.ReadAll(resp.Body)
	fmt.Println(string(content))
}

type Person struct {
	Name string
	Age  int
}

// Post
func Test2() {
	person := Person{
		Name: "zs",
		Age:  20,
	}
	json, _ := json.Marshal(person)
	reader := bytes.NewReader(json)

	resp, err := http.Post("http://www.baidu.com", "application/json", reader)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

// 客户端
// http.Client{}结构体，可提供的配置项总共有四个:
// Transport:配置 Http 客户端数据传输相关的配置项，没有就采用默认的策略
// Timeout：请求超时时间配置
// Jar：Cookie 相关配置
// CheckRedirect：重定向配置
func Test3() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}

func Test4() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	req.Header.Add("Authorization", "123")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}

// 服务端
func Test5() {
	//简单版
	http.ListenAndServe("localhost:8080", nil)

	//自定义
	server := &http.Server{
		Addr:              "localhost:8080",
		Handler:           nil,
		TLSConfig:         nil,
		ReadTimeout:       0,
		WriteTimeout:      0,
		ReadHeaderTimeout: 0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil}
	server.SetKeepAlivesEnabled(false)
}

// 路由
// 首先自定义一个结构体实现Handler接口中的ServeHTTP(ResponseWriter, *Request)方法，再调用http.handle()函数即可
type MyHandler struct{}

func (*MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "hello world")
}
func Test6() {
	http.Handle("/index", &MyHandler{})
	http.ListenAndServe("localhost:8080", nil)
}

// 可以直接http.handlerFunc函数，我们只需要写处理函数即可，从而不用创建结构体。
// 其内部是使用了适配器类型HandlerFunc,HandlerFunc 类型是一个适配器，允许将普通函数用作 HTTP 的处理器。
// 如果 f 是具有适当签名的函数，HandlerFunc(f)是调用 f 的 Handler
func Test7() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "hello world")
	})
	http.ListenAndServe("localhost:8080", nil)
}

// 反向代理
func Test8() {
	http.HandleFunc("/forward", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = "www.baidu.com"
			req.URL.Path = "/index"
		}
		proxy := httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})
	http.ListenAndServe("localhost:8080", nil)
}
