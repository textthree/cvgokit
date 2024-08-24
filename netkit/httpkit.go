package netkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

// 发送 Http Get 请求
// args[0] 要发送的数据
func HttpGet(url string, args ...map[string]string) string {
	var queryStringBuffer bytes.Buffer
	if len(args) > 0 {
		counter := 0
		length := len(args[0])
		for k, v := range args[0] {
			if counter == 0 {
				queryStringBuffer.WriteString("?")
			}
			queryStringBuffer.WriteString(k)
			queryStringBuffer.WriteString("=")
			queryStringBuffer.WriteString(v)
			counter++
			length--
			if length > 0 {
				queryStringBuffer.WriteString("&")
			}
		}
	}
	queryString := queryStringBuffer.String()
	//Trace(url + queryString)
	// 超时时间：6 秒
	client := &http.Client{Timeout: 6 * time.Second}
	resp, err := client.Get(url + queryString)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	return result.String()
}

// 发送 Http Post 请求
// args[0] 要发送的数据(json字符串)
// args[1] 头信息
func HttpPostJson(url string, args ...[]byte) (res []byte, error error) {
	var params []byte
	if len(args) > 0 {
		params = args[0]
	}
	postData := strings.NewReader(string(params))
	// 6 秒超时
	client := &http.Client{Timeout: 6 * time.Second}
	request, err := http.NewRequest("POST", url, postData)
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	if err != nil {
		error = err
	}
	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		error = err
		return
	}
	defer response.Body.Close()
	res, err = ioutil.ReadAll(response.Body)
	return
}

// 发送 Http Post 请求
// data 要发送的数据
// header 头信息 如: headers := map[string]string{"token":"123"}
func Post(url string, data, headers map[string]string) (res []byte, error error) {
	var params []byte
	params, error = json.Marshal(data)
	if error != nil {
		fmt.Println("httpkit Post error")
	}
	postData := strings.NewReader(string(params))
	// 6 秒超时
	client := &http.Client{Timeout: 6 * time.Second}
	request, err := http.NewRequest("POST", url, postData)
	if err != nil {
		error = err
	}
	// 添加头信息
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	for key, header := range headers {
		request.Header.Set(key, header)
	}
	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		error = err
		return
	}
	defer response.Body.Close()
	res, err = ioutil.ReadAll(response.Body)
	return
}

// 检查是否有网络
func NetWorkStatus() bool {
	cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
	fmt.Println("NetWorkStatus Start:", time.Now().Unix())
	err := cmd.Run()
	fmt.Println("NetWorkStatus End  :", time.Now().Unix())
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		fmt.Println("Net Status , OK")
	}
	return true
}

// 使用代理发送 get 请求
func HttpGetWithProxy(proxyServer string) {
	// 代理服务器地址
	proxyStr := "http://proxy.example.com:8080"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		panic(err)
	}

	// 创建自定义的 http.Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 创建自定义的 http.Client
	client := &http.Client{
		Transport: transport,
	}

	// 目标 URL
	targetURL := "http://example.com"

	// 使用自定义的 http.Client 发送请求
	response, err := client.Get(targetURL)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// 输出响应体
	fmt.Println(string(body))
}
