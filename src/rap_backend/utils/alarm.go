package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const ALARM_URL = "http://base-service.airudder.com/alarmCommon"
const ALARM_BOT = "https://open.feishu.cn/open-apis/bot/v2/hook/7dc6f500-71fa-4992-b558-0be59dce7d09"

var Env string

func AlarmHttp(title, msg string) error {
	var alarm struct {
		Title    string `json:"title" binding:"required"`
		Msg      string `json:"msg" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	alarm.Title = title
	alarm.Msg = msg
	alarm.Category = "rap"

	var result struct {
		Code    int    `json:"code" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := DoHttpPostJson(ALARM_URL, &alarm, &result); err != nil {
		return err
	}
	if result.Code != 0 {
		return errors.New(result.Message)
	}
	return nil
}

type AlarmBotReq struct {
	Title    string     `json:"title"`
	Msg      string     `json:"msg"`
	Category string     `json:"category"`
	URL      string     `json:"url"`
	MsgType  string     `json:"msg_type"`
	Content  BotContent `json:"content"`
}

type BotContent struct {
	Text string `json:"text"`
}

// Alarm ...
func AlarmWithBot(title, body string) error {
	msg := `【报警】- %s 
标题：%s
内容：%s`
	reqBody := AlarmBotReq{
		MsgType: "text",
		Content: BotContent{
			Text: fmt.Sprintf(msg, Env, title, body),
		},
	}

	var result struct {
		Code    int    `json:"code" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := DoHttpPostJson(ALARM_BOT, reqBody, &result); err != nil {
		return err
	}
	if result.Code != 0 {
		return errors.New(result.Message)
	}
	return nil
}

// Alarm ...
func AlarmWithHttp(url, name, supervisor, title, body, dept, cc string, level int) {

	if Env == "dev" {
		level = 3
		dept = "交互2.0测试报警群"
	}
	name = "junliang.chen"
	supervisor = "junliang.chen"
	req := Req{
		Name:       name,
		Supervisor: supervisor,
		Title:      title,
		Dept:       dept,
		Body:       body,
		Level:      level,
		Cc:         cc,
	}
	data, _ := json.Marshal(req)
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {

	} else {
		// 自行处理
	}
}
