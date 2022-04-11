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
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/sirupsen/logrus"

	"github.com/AlecAivazis/survey/v2"
)

func NewSession(cookie string, num int) *Session {
	// TODO 增加并发请求
	if num < 1 {
		num = 1
	}
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
		client: client,
		num:    num,
	}
}

type Session struct {
	client *resty.Client
	num    int

	UserID       string
	Address      *AddressItem
	BarkId       string
	Cart         Cart
	Order        Order
	PackageOrder PackageOrder
	PayType      int
	CartMode     int
}

func (s *Session) Clone() *Session {
	return &Session{
		client:   s.client,
		UserID:   s.UserID,
		Address:  s.Address,
		BarkId:   s.BarkId,
		PayType:  s.PayType,
		CartMode: s.CartMode,
	}
}

func (s *Session) buildHeader() http.Header {
	header := make(http.Header)
	header.Set("ddmc-city-number", s.Address.CityNumber)
	header.Set("ddmc-os-version", "undefined")
	header.Set("ddmc-channel", "applet")
	header.Set("ddmc-build-version", "2.82.0")
	header.Set("ddmc-app-client-id", "4")
	header.Set("ddmc-ip", "")
	header.Set("ddmc-api-version", "9.49.2")
	header.Set("ddmc-station-id", s.Address.StationId)
	header.Set("ddmc-uid", s.UserID)
	return header
}

func (s *Session) buildURLParams(needAddress bool) url.Values {
	params := url.Values{}
	params.Add("api_version", "9.49.0")
	params.Add("app_version", "2.81.0")
	params.Add("applet_source", "")
	params.Add("app_client_id", "3")
	params.Add("h5_source", "")
	params.Add("sharer_uid", "")
	params.Add("s_id", "")
	params.Add("openid", "")
	if needAddress {
		params.Add("station_id", s.Address.StationId)
		params.Add("city_number", s.Address.CityNumber)
	}
	return params
}

func (s *Session) Choose() error {
	if err := s.chooseAddr(); err != nil {
		return err
	}
	if err := s.choosePay(); err != nil {
		return err
	}
	if err := s.chooseCartMode(); err != nil {
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
	for _, v := range addrMap {
		addrs = append(addrs, v.Location.Address)
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
	logrus.Infof("已选择收货地址: %s", s.Address.Location.Address)
	return nil
}

const (
	paymentAlipay = "支付宝"
	paymentWechat = "微信"
)

func (s *Session) choosePay() error {
	var payName string
	sv := &survey.Select{
		Message: "请选择支付方式",
		Options: []string{paymentWechat, paymentAlipay},
		Default: paymentWechat,
	}
	if err := survey.AskOne(sv, &payName); err != nil {
		return fmt.Errorf("选择支付方式错误: %v", err)
	}

	switch payName {
	case paymentAlipay:
		s.PayType = 2
	case paymentWechat:
		s.PayType = 4
	default:
		return fmt.Errorf("无法识别的支付方式: %s", payName)
	}
	return nil
}

const (
	cartModeAvailable = "结算所有有效商品(不包括换购)"
	cartModeAll       = "结算所有勾选商品(包括换购)"
)

func (s *Session) chooseCartMode() error {
	var cartDesc string
	sv := &survey.Select{
		Message: "请选择购物车商品结算模式",
		Options: []string{cartModeAvailable, cartModeAll},
		Default: cartModeAvailable,
	}
	if err := survey.AskOne(sv, &cartDesc); err != nil {
		return fmt.Errorf("选择购物车商品结算模式错误: %v", err)
	}

	switch cartDesc {
	case cartModeAvailable:
		s.CartMode = 1
	case cartModeAll:
		s.CartMode = 2
	default:
		return fmt.Errorf("无法识别的购物车商品结算模式: %s", cartDesc)
	}
	return nil
}
