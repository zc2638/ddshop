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

package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func NewSession(cookie string, interval int64) *Session {
	if !strings.HasPrefix(cookie, "DDXQSESSID=") {
		cookie = "DDXQSESSID=" + cookie
	}

	header := make(http.Header)
	header.Set("Host", "maicai.api.ddxq.mobi")
	header.Set("user-agent", "Mozilla/5.0 (Linux; Android 9; LIO-AN00 Build/LIO-AN00; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/92.0.4515.131 Mobile Safari/537.36 xzone/9.47.0 station_id/null")
	header.Set("accept", "application/json, text/plain, */*")
	header.Set("content-type", "application/x-www-form-urlencoded")
	header.Set("origin", "https://wx.m.ddxq.mobi")
	header.Set("x-requested-with", "com.yaya.zone")
	header.Set("sec-fetch-site", "same-site")
	header.Set("sec-fetch-mode", "cors")
	header.Set("sec-fetch-dest", "empty")
	header.Set("referer", "https://wx.m.ddxq.mobi/")
	header.Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	header.Set("cookie", cookie)

	client := resty.New()
	client.Header = header
	return &Session{
		client:   client,
		interval: interval,

		apiVersion:  "9.49.2",
		appVersion:  "2.82.0",
		channel:     "applet",
		appClientID: "4",
	}
}

type Session struct {
	client   *resty.Client
	interval int64 // 间隔请求时间(ms)

	channel     string
	apiVersion  string
	appVersion  string
	appClientID string

	UserID  string
	PayType int64
	Address *AddressItem
	Reserve ReserveTime
}

func (s *Session) Clone() *Session {
	return &Session{
		client:   s.client,
		interval: s.interval,

		channel:     s.channel,
		apiVersion:  s.apiVersion,
		appVersion:  s.appVersion,
		appClientID: s.appClientID,

		UserID:  s.UserID,
		Address: s.Address,
		PayType: s.PayType,
		Reserve: s.Reserve,
	}
}

func (s *Session) execute(ctx context.Context, request *resty.Request, method, url string) (*resty.Response, error) {
	if ctx != nil {
		request.SetContext(ctx)
	}
	resp, err := request.Execute(method, url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("statusCode: %d, body: %s", resp.StatusCode(), resp.String())
	}

	result := gjson.ParseBytes(resp.Body())
	code := result.Get("code").Num
	switch code {
	case 0:
		return resp, nil
	case -3000, -3001:
		logrus.Warningf("当前人多拥挤(%v): %s", code, resp.String())
	default:
		return nil, fmt.Errorf("无法识别的状态码: %v", resp.String())
	}
	logrus.Warningf("将在 %dms 后重试", s.interval)
	time.Sleep(time.Duration(s.interval) * time.Millisecond)
	return s.execute(nil, request, method, url)
}

func (s *Session) buildHeader() http.Header {
	header := make(http.Header)
	header.Set("ddmc-city-number", s.Address.CityNumber)
	header.Set("ddmc-os-version", "undefined")
	header.Set("ddmc-channel", s.channel)
	header.Set("ddmc-api-version", s.apiVersion)
	header.Set("ddmc-build-version", s.appVersion)
	header.Set("ddmc-app-client-id", s.appClientID)
	header.Set("ddmc-ip", "")
	header.Set("ddmc-station-id", s.Address.StationId)
	header.Set("ddmc-uid", s.UserID)
	if len(s.Address.Location.Location) == 2 {
		header.Set("ddmc-longitude", strconv.FormatFloat(s.Address.Location.Location[0], 'f', -1, 64))
		header.Set("ddmc-latitude", strconv.FormatFloat(s.Address.Location.Location[1], 'f', -1, 64))
	}
	return header
}

func (s *Session) buildURLParams(needAddress bool) url.Values {
	params := url.Values{}
	params.Add("channel", s.channel)
	params.Add("api_version", s.apiVersion)
	params.Add("app_version", s.appVersion)
	params.Add("app_client_id", s.appClientID)
	params.Add("applet_source", "")
	params.Add("h5_source", "")
	params.Add("sharer_uid", "")
	params.Add("s_id", "")
	params.Add("openid", "")

	params.Add("uid", s.UserID)
	if needAddress {
		params.Add("address_id", s.Address.Id)
		params.Add("station_id", s.Address.StationId)
		params.Add("city_number", s.Address.CityNumber)
		if len(s.Address.Location.Location) == 2 {
			params.Add("longitude", strconv.FormatFloat(s.Address.Location.Location[0], 'f', -1, 64))
			params.Add("latitude", strconv.FormatFloat(s.Address.Location.Location[1], 'f', -1, 64))
		}
	}

	params.Add("device_token", "")
	params.Add("nars", "")
	params.Add("sesi", "")
	return params
}

func (s *Session) SetReserve(reserve ReserveTime) {
	s.Reserve = reserve
}

func (s *Session) Choose(payType string) error {
	if err := s.chooseAddr(); err != nil {
		return err
	}
	if err := s.choosePay(payType); err != nil {
		return err
	}
	return nil
}

func (s *Session) chooseAddr() error {
	addrMap, err := s.GetAddress()
	if err != nil {
		return fmt.Errorf("获取收货地址失败: %v", err)
	}
	addrs := make([]string, 0, len(addrMap))
	for k := range addrMap {
		addrs = append(addrs, k)
	}

	if len(addrs) == 1 {
		address := addrMap[addrs[0]]
		s.Address = &address
		logrus.Infof("默认收货地址: %s %s", s.Address.Location.Address, s.Address.AddrDetail)
		return nil
	}

	var addr string
	sv := &survey.Select{
		Message: "请选择收货地址",
		Options: addrs,
	}
	if err := survey.AskOne(sv, &addr); err != nil {
		return fmt.Errorf("选择收货地址错误: %v", err)
	}

	address, ok := addrMap[addr]
	if !ok {
		return errors.New("请选择正确的收货地址")
	}
	s.Address = &address
	logrus.Infof("已选择收货地址: %s %s", s.Address.Location.Address, s.Address.AddrDetail)
	return nil
}

const (
	PaymentAlipay    = "alipay"
	PaymentAlipayStr = "支付宝"
	PaymentWechat    = "wechat"
	PaymentWechatStr = "微信"
)

func (s *Session) choosePay(payType string) error {
	if payType == "" {
		sv := &survey.Select{
			Message: "请选择支付方式",
			Options: []string{PaymentWechatStr, PaymentAlipayStr},
			Default: PaymentWechatStr,
		}
		if err := survey.AskOne(sv, &payType); err != nil {
			return fmt.Errorf("选择支付方式错误: %v", err)
		}
	}

	// 2支付宝，4微信，6小程序支付
	switch payType {
	case PaymentAlipay, PaymentAlipayStr:
		s.PayType = 2
		logrus.Info("已选择支付方式：支付宝")
	case PaymentWechat, PaymentWechatStr:
		s.PayType = 4
		logrus.Info("已选择支付方式：微信")
	default:
		return fmt.Errorf("无法识别的支付方式: %s", payType)
	}
	return nil
}
