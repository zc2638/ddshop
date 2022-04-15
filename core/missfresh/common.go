package missfresh

import "strings"

const (
	HOST = "https://as-vip.missfresh.cn"

	User_API           = "/as/user/address/default/view"
	REGION_API         = "/as/miss/forerunner/location"   // 1 获取region，请求头需要使用
	USERINFO_API       = "/v1/auth/info"                  // 2
	ADDRESS_API        = "/v1/customer/address/list"      // 3
	CART_PROMOTION_API = "/as/portal/cart/sync/promotion" // 4 获取购物车
	ORDER_API          = "/as/portal/order/createOrder"   // 5 创建订单
)

const maxRetryCount = 1

func (s *Session) CreateBodyParam() map[string]interface{} {
	m := map[string]interface{}{
		"common": map[string]interface{}{
			"accessToken": strings.ReplaceAll(s.cfg.Token, "=", "%3D"),
			//"accessToken":    s.cfg.Token,
			"fromSource":     "share_miniprogram_80420239",
			"retailType":     "",
			"sourceDeviceId": s.cfg.DeviceID,
			"deviceId":       s.cfg.DeviceID,
			"deviceCenterId": s.cfg.DeviceCenterID,
			"env":            "weixin_app",
			"platform":       "weixin_app",
			"model":          "MacBookPro18%2C4",
			"screenHeight":   736,
			"screenWidth":    414,
			"version":        "9.9.92.11",
			"addressCode":    "",
			"stationCode":    "",
			"bigWarehouse":   "",
			"deliveryType":   0,
			"chromeType":     0,
			"currentLng":     "",
			"currentLat":     "",
			"sellerId":       0,
			"mfplatform":     "weixin_app",
			"mfenv":          "wxapp",
		},
		"param": map[string]interface{}{},
	}

	return m
}
