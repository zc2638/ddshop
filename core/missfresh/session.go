package missfresh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/go-resty/resty/v2"
	"github.com/zc2638/ddshop/pkg/notice"
)

func NewSession(cfg *Config, noticeIns notice.Interface) (*Session, error) {

	header := make(http.Header)
	header.Set("Host", strings.TrimPrefix(HOST, "https"))
	header.Set("content-type", "application/json")
	header.Set("accept", "*/*")
	header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 11_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E217 MicroMessenger/6.8.0(0x16080000) NetType/WIFI Language/en Branch/Br_trunk MiniProgramEnv/Mac")
	header.Set("referer", "https://servicewechat.com/wxebf773691904eee9/985/page-frame.html")
	header.Set("accept-language", "zh-cn")
	header.Set("platform", "weixin_app")
	header.Set("accept-encoding", "gzip, deflate, br")
	header.Set("accesstoken", cfg.Token)
	header.Set("request-id", strings.ReplaceAll(cfg.DeviceID, "-", ""))

	client := resty.New().SetBaseURL(HOST)
	client.Header = header

	sess := &Session{
		cfg:       cfg,
		noticeIns: noticeIns,
		client:    client,
		successCh: make(chan struct{}, 1),
		stopCh:    make(chan struct{}, 1),

		channel: "applet",
	}

	if err := sess.GetDefaultAddressView(context.Background()); err != nil {
		return nil, fmt.Errorf("获取用户view失败：%v", err)
	}

	if err := sess.GetRegion(context.Background()); err != nil {
		return nil, fmt.Errorf("获取地区信息失败：%v", err)
	}
	region := `{"address_code":"%s","station_code":"%s","delivery_type":%d,"bigWarehouse":"%s","type":%d`
	header.Set("x-region", fmt.Sprintf(region,
		sess.RegionInfo.AreaCode,
		sess.RegionInfo.StationCode,
		sess.RegionInfo.DeliveryType,
		sess.RegionInfo.BigWarehouse,
		sess.RegionInfo.Type))

	if err := sess.Choose(); err != nil {
		return nil, err
	}
	return sess, nil
}

type Session struct {
	cfg       *Config
	noticeIns notice.Interface
	client    *resty.Client
	successCh chan struct{}
	stopCh    chan struct{}

	channel string

	DefaultAddressInfo *DefaultAddressInfo
	PayType            int64
	Address            *AddressItem
	RegionInfo         *RegionInfo
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
	addrMap, err := s.GetAddress(context.Background())
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
		logrus.Infof("默认收货地址: %s", s.Address.AddressDetail)
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
	logrus.Infof("已选择收货地址: %s", s.Address.AddressDetail)
	return nil
}

const (
	PaymentAlipay         = "alipay"
	PaymentAlipayStr      = "支付宝"
	PaymentWechat         = "wechat"
	PaymentWechatStr      = "微信"
	PaymentStoredValue    = "stored"
	PaymentStoredValueStr = "储值"
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

	// 2支付宝，4微信，6储值
	switch payType {
	case PaymentAlipay, PaymentAlipayStr:
		s.PayType = 2
		logrus.Info("已选择支付方式：支付宝")
	case PaymentWechat, PaymentWechatStr:
		s.PayType = 4
		logrus.Info("已选择支付方式：微信")
	case PaymentStoredValue, PaymentStoredValueStr:
		s.PayType = 6
		logrus.Info("已选择支付方式：储值")
	default:
		return fmt.Errorf("无法识别的支付方式: %s", payType)
	}
	return nil
}

func (s *Session) execute(ctx context.Context, request *resty.Request, method, url string, count int) (*resty.Response, error) {
	if ctx != nil {
		request.SetContext(ctx)
	}

	request.SetHeader("Host", strings.TrimPrefix(HOST, "https://"))
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
	case 30007:
		return nil, fmt.Errorf("状态码：%v，错误：%v", code, resp.String())
	default:
		if count <= 0 {
			return nil, fmt.Errorf("无法识别的状态码: %v", resp.String())
		}
		logrus.Warningf("尝试次数: %d, 无法识别的状态码: %v", count, resp.String())
	}
	count--
	return s.execute(nil, request, method, url, count)
}

func (s *Session) buildURLParams(needAddress bool) url.Values {
	params := url.Values{}

	params.Add("fromSource", "share_miniprogram_80420239")
	params.Add("version", "9.9.92.11")
	params.Add("platform", "weixin_app")
	params.Add("device_id", s.cfg.DeviceID)
	params.Add("deviceCenterId", s.cfg.DeviceCenterID)
	params.Add("business_type", "")
	params.Add("mfplatform", "weixin_app")
	params.Add("mfenv", "wxapp")
	params.Add("sellerId", "0")

	if needAddress {
		params.Add("sellerId", strconv.Itoa(s.RegionInfo.SellerId))
		params.Add("currentLng", s.RegionInfo.Lng)
		params.Add("currentLat", s.RegionInfo.Lat)
	}

	return params
}

func (s *Session) Run(ctx context.Context) error {
	err := s.run(ctx)
	if err != nil {
		switch err {
		//case ErrorNoValidProduct:
		//	sleepInterval := 30
		//	logrus.Errorf("购物车中无有效商品，请先前往app添加或勾选，%d 秒后重试！", sleepInterval)
		//	time.Sleep(time.Duration(sleepInterval) * time.Second)
		//case ErrorNoReserveTime:
		//	sleepInterval := 3 + rand.Intn(6)
		//	logrus.Warningf("暂无可预约的时间，%d 秒后重试！", sleepInterval)
		//	time.Sleep(time.Duration(sleepInterval) * time.Second)
		default:
			logrus.Error(err)
		}
	}
	return err
}

func (s *Session) run(ctx context.Context) error {
	logrus.Info("=====> 获取购物车中有效商品")

	if err := s.GetCart(ctx); err != nil {
		return fmt.Errorf("获取购物车商品失败: %v", err)
	}

	return nil
	//products := cartData["products"].([]map[string]interface{})
	//for k, v := range products {
	//	logrus.Infof("[%v] %s 数量：%v 总价：%s", k, v["product_name"], v["count"], v["total_price"])
	//}
	//
	//for {
	//
	//	logrus.Info("=====> 生成订单信息")
	//	err := s.CreateOrder(ctx)
	//	if err != nil {
	//		return fmt.Errorf("检查订单失败: %v", err)
	//	}
	//	logrus.Infof("订单总金额：%v\n", checkOrderData["price"])
	//
	//	logrus.Infof("=====> 提交订单中, 预约时间段(%s)", timeRange)
	//	if err := sess.CreateOrder(context.Background()); err != nil {
	//		logrus.Errorf("提交订单(%s)失败: %v", err)
	//		return err
	//	}
	//
	//	if err := s.noticeIns.Notice("抢菜成功", "叮咚买菜 抢菜成功，请尽快支付！"); err != nil {
	//		logrus.Warningf("通知失败: %v", err)
	//	}
	//	return nil
	//}
}
