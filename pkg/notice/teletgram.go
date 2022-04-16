package notice

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type TelegramBotConf struct {
	Token   string `json:"token"`
	ChatID  string `json:"chat_id"`
	Content string `json:"content"`
	Title   string `json:"title"`
}

type TelegramBot struct {
	Conf TelegramBotConf
}

func NewTelegramBot(cf *TelegramBotConf) *TelegramBot {
	return &TelegramBot{
		Conf: *cf,
	}
}

func (b *TelegramBot) Name() string {
	return "Telegram"
}

func (b *TelegramBot) Send(title, content string) error {
	sendUrl := fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", b.Conf.Token)
	v := url.Values{}
	v.Add("chat_id", b.Conf.ChatID)
	v.Add("text", fmt.Sprintf("[%v]%v", title, content))
	payload := strings.NewReader(v.Encode())
	resp, err := http.Post(sendUrl, "application/x-www-form-urlencoded", payload)
	if err != nil {
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode: %d, body: %v", resp.StatusCode, string(result))
	}
	return nil
}
