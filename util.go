package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

func DoGet(oriUrl string, params map[string]string) (int, string, error) {
	// 添加查询参数
	u, _ := url.Parse(oriUrl)
	q := u.Query()
	if params != nil {
		for key, value := range params {
			q.Add(key, value)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Println(err)
		return 0, "", err
	}
	req.AddCookie(&http.Cookie{
		Name:  "zentaosid",
		Value: sessionId,
	})

	// Create a new client
	client := &http.Client{
		Jar: cookies,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, "", err
	}
	defer resp.Body.Close()
	log.Println("Get status:", resp.Status, oriUrl)
	if resp.StatusCode != 200 {
		return 0, "", errors.New("Get error with " + resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(content), err
}

// DoPost 发送 POST 请求，参数以表单形式传递
func DoPost(oriUrl string, params map[string]string) (int, string, error) {
	// 将参数编码为表单数据
	formData := url.Values{}
	if params != nil {
		for key, value := range params {
			formData.Add(key, value)
		}
	}

	// 创建请求
	req, err := http.NewRequest("POST", oriUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Println(err)
		return 0, "", err
	}

	// 设置请求头，指定内容类型为表单数据
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 添加 Cookie
	req.AddCookie(&http.Cookie{
		Name:  "zentaosid",
		Value: sessionId,
	})

	// 创建 HTTP 客户端
	client := &http.Client{
		Jar: cookies,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return 0, "", err
	}
	defer resp.Body.Close()

	// 日志记录响应状态
	log.Println("Post status:", resp.Status, oriUrl)

	// 检查响应状态码
	if resp.StatusCode != 200 {
		return 0, "", errors.New("Post error with " + resp.Status)
	}

	// 读取响应内容
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0, "", err
	}

	// 返回状态码、响应内容和错误
	return resp.StatusCode, string(content), nil
}

// 生成32位uuid（伪实现，真实可用uuid库）
func genZentaoSID() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func getThisWeekMonday() time.Time {
	// 获取当前时间
	now := time.Now()

	// 计算本周一的日期
	year, month, day := now.Date()
	weekDay := now.Weekday()
	// 如果今天是周一，Weekday() 返回 time.Monday
	// 如果今天是周二，Weekday() 返回 time.Tuesday，依此类推
	// time.Monday 的值为 0，time.Tuesday 的值为 1，依此类推
	// 因此，本周一的日期为当前日期减去 (今天是周几 - 0) 天
	monday := time.Date(year, month, day-int(weekDay)-int(time.Monday)+1, 0, 0, 0, 0, now.Location())

	return monday
}
