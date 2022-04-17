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
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
)

const barkURL = "https://api.day.app/push"

type BarkConfig struct {
	Server string `json:"server"`
	Key    string `json:"key"`
}

func NewBark(cfg *BarkConfig) Engine {
	if cfg.Key == "" {
		return nil
	}
	return &bark{cfg: cfg}
}

type bark struct {
	cfg *BarkConfig
}

func (b *bark) Name() string {
	return "Bark"
}

func (b *bark) Send(title, body string) error {
	data := &barkData{
		DeviceKey: b.cfg.Key,
		Title:     title,
		Body:      body,
		Sound:     "alarm.caf",
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}

	uri := barkURL
	if b.cfg.Server != "" {
		uri = path.Join(b.cfg.Server, "push")
	}
	resp, err := http.Post(uri, "application/json; charset=utf-8", bytes.NewReader(bs))
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

type barkData struct {
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	Badge     int    `json:"badge,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	Url       string `json:"url,omitempty"`
}
