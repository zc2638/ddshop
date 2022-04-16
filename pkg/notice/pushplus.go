// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package notice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const pushPlusURL = "http://www.pushplus.plus/send"

type PushPlusConfig struct {
	Token string `json:"token"`
}

func NewPushPlus(cfg *PushPlusConfig) Engine {
	if cfg.Token == "" {
		return nil
	}
	return &pushPlus{token: cfg.Token}
}

type pushPlusResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type pushPlus struct {
	token string
}

func (p *pushPlus) Name() string {
	return "pushplus"
}

func (p *pushPlus) Send(title, body string) error {
	data := map[string]string{
		"token":   p.token,
		"title":   title,
		"content": body,
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := http.Post(pushPlusURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}

	ba, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("statusCode: %d, body: %v", resp.StatusCode, string(ba))
	}

	var res pushPlusResult
	if err := json.Unmarshal(ba, &res); err != nil {
		return err
	}
	if res.Code != 200 {
		return errors.New(res.Msg)
	}
	return nil
}
