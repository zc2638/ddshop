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

package ddmc

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/zc2638/ddshop/pkg/notice"
)

var (
	ErrorNoValidProduct     = errors.New("无有效商品")
	ErrorNoReserveTime      = errors.New("无可预约时间段")
	ErrorOutStock           = errors.New("部分商品已缺货")
	ErrorProductChange      = errors.New("商品信息发生变化")
	ErrorReserveTimeExpired = errors.New("送达时间已失效")
)

func NewSession(cfg *Config, noticeIns notice.Interface) (*Session, error) {
	cookie := cfg.Cookie
	if !strings.HasPrefix(cookie, "DDXQSESSID=") {
		cookie = "DDXQSESSID=" + cookie
	}

	header := make(http.Header)
	header.Set("Host", "maicai.api.ddxq.mobi")
	header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E217 MicroMessenger/6.8.0(0x16080000) NetType/WIFI Language/en Branch/Br_trunk MiniProgramEnv/Mac")
	header.Set("accept", "application/json, text/plain, */*")
	header.Set("content-type", "application/x-www-form-urlencoded")
	header.Set("origin", "https://wx.m.ddxq.mobi")
	header.Set("sec-fetch-site", "same-site")
	header.Set("sec-fetch-mode", "cors")
	header.Set("sec-fetch-dest", "empty")
	header.Set("referer", "https://wx.m.ddxq.mobi/")
	header.Set("accept-language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	header.Set("cookie", cookie)

	client := resty.New()
	client.Header = header

	if cfg.Channel != 3 && cfg.Channel != 4 {
		cfg.Channel = 4
	}
	sess := &Session{
		cfg:       cfg,
		noticeIns: noticeIns,
		client:    client,

		apiVersion:  "9.49.2",
		appVersion:  "2.82.0",
		channel:     "applet",
		appClientID: strconv.Itoa(cfg.Channel),
	}

	if err := sess.GetUser(); err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	if err := sess.Choose(); err != nil {
		return nil, err
	}
	return sess, nil
}

type Session struct {
	cfg       *Config
	noticeIns notice.Interface
	client    *resty.Client

	channel     string
	apiVersion  string
	appVersion  string
	appClientID string

	UserID  string
	PayType int64
	Address *AddressItem

	cartData         map[string]interface{}
	multiReserveTime []ReserveTime
	checkOrderData   map[string]interface{}
}

func (s *Session) Run(ctx context.Context) error {
	err := s.run(ctx)
	if err != nil {
		switch err {
		case ErrorNoValidProduct:
			sleepInterval := 30
			logrus.Errorf("购物车中无有效商品，请先前往app添加或勾选，%d 秒后重试！", sleepInterval)
			time.Sleep(time.Duration(sleepInterval) * time.Second)
		case ErrorNoReserveTime:
			sleepInterval := 3 + rand.Intn(6)
			logrus.Warningf("暂无可预约的时间，%d 秒后重试！", sleepInterval)
			time.Sleep(time.Duration(sleepInterval) * time.Second)
		default:
			logrus.Error(err)
		}
	}
	return err
}

func (s *Session) run(ctx context.Context) error {
	if s.cartData == nil {
		logrus.Info("=====> 获取购物车中有效商品")
		if err := s.CartAllCheck(ctx); err != nil {
			return fmt.Errorf("全选购物车商品失败: %v", err)
		}
		cartData, err := s.GetCart(ctx)
		if err != nil {
			return err
		}
		s.cartData = cartData
	}

	products := s.cartData["products"].([]map[string]interface{})
	for k, v := range products {
		logrus.Infof("[%v] %s 数量：%v 总价：%s", k, v["product_name"], v["count"], v["total_price"])
	}

	if len(s.multiReserveTime) == 0 {
		logrus.Info("=====> 获取可预约时间")
		multiReserveTime, err := s.GetMultiReserveTime(ctx, products)
		if err != nil {
			return fmt.Errorf("获取可预约时间失败: %v", err)
		}
		if len(multiReserveTime) == 0 {
			return ErrorNoReserveTime
		}
		logrus.Infof("发现可用的配送时段!")
		s.multiReserveTime = multiReserveTime
	}
	reserveTime := s.multiReserveTime[0]

	if s.checkOrderData == nil {
		logrus.Info("=====> 生成订单信息")
		checkOrderData, err := s.CheckOrder(ctx, s.cartData, &reserveTime)
		if err != nil {
			return fmt.Errorf("检查订单失败: %v", err)
		}
		s.checkOrderData = checkOrderData
		logrus.Infof("订单总金额：%v\n", s.checkOrderData["price"])
	}

	startTime := time.Unix(int64(reserveTime.StartTimestamp), 0).Format("2006/01/02 15:04:05")
	endTime := time.Unix(int64(reserveTime.EndTimestamp), 0).Format("2006/01/02 15:04:05")
	timeRange := startTime + "——" + endTime

	logrus.Infof("=====> 提交订单中, 预约时间段(%s)", timeRange)
	if err := s.CreateOrder(context.Background(), s.cartData, s.checkOrderData, &reserveTime); err != nil {
		logrus.Errorf("提交订单(%s)失败: %v", timeRange, err)
		return err
	}

	if err := s.noticeIns.Notice("抢菜成功", "叮咚买菜 抢菜成功，请尽快支付！"); err != nil {
		logrus.Warningf("通知失败: %v", err)
	}
	return nil
}

func (s *Session) execute(ctx context.Context, request *resty.Request, method, url string, count int) (*resty.Response, error) {
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
	case -3000:
		msg := result.Get("msg").Str
		if count == 0 {
			return nil, fmt.Errorf("当前人多拥挤(%v): %s", code, msg)
		}
		logrus.Warningf("将在 %dms 后重试, 当前人多拥挤(%v): %s", s.cfg.Interval, code, msg)
	case -3001, -3100: // -3001创建订单, -3100检查订单
		msg := result.Get("tips.limitMsg").Str
		if count == 0 {
			return nil, fmt.Errorf("当前拥挤(%v): %s", code, msg)
		}

		interval := int64(result.Get("tips.duration").Num)
		if interval == 0 {
			interval = s.cfg.Interval
		}
		logrus.Warningf("将在 %dms 后重试, 当前人多拥挤(%v): %s", interval, code, msg)
		time.Sleep(time.Duration(interval) * time.Millisecond)
	case 5001:
		s.cartData = nil
		s.checkOrderData = nil
		return nil, ErrorOutStock
	case 5003:
		s.cartData = nil
		s.checkOrderData = nil
		return nil, ErrorProductChange
	case 5004:
		s.multiReserveTime = nil
		return nil, ErrorReserveTimeExpired
	default:
		s.checkOrderData = nil
		return nil, fmt.Errorf("无法识别的状态码: %v", resp.String())
	}
	count--
	return s.execute(nil, request, method, url, count)
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

func (s *Session) Choose() error {
	if err := s.chooseAddr(); err != nil {
		return err
	}
	if err := s.choosePay(); err != nil {
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

func (s *Session) choosePay() error {
	payType := s.cfg.PayType
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
