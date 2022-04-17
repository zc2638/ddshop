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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

func (s *Session) CheckOrder(ctx context.Context, cartData map[string]interface{}, reserveTime *ReserveTime) (map[string]interface{}, error) {
	urlPath := "https://maicai.api.ddxq.mobi/order/checkOrder"

	packagesInfo := make(map[string]interface{})
	for k, v := range cartData {
		packagesInfo[k] = v
	}
	packagesInfo["reserved_time"] = map[string]interface{}{
		"reserved_time_start": reserveTime.StartTimestamp,
		"reserved_time_end":   reserveTime.EndTimestamp,
	}
	packagesJson, err := json.Marshal([]interface{}{packagesInfo})
	if err != nil {
		return nil, fmt.Errorf("marshal products info failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Add("packages", string(packagesJson))
	params.Add("user_ticket_id", "default")
	params.Add("freight_ticket_id", "default")
	params.Add("is_use_point", "0")
	params.Add("is_use_balance", "0")
	params.Add("is_buy_vip", "0")
	params.Add("coupons_id", "")
	params.Add("is_buy_coupons", "0")
	params.Add("check_order_type", "0")
	params.Add("is_support_merge_payment", "0")
	params.Add("showData", "true")
	params.Add("showMsg", "false")

	req := s.client.R()
	req.Header = s.buildHeader()
	req.SetBody(strings.NewReader(params.Encode()))
	resp, err := s.execute(ctx, req, http.MethodPost, urlPath, s.cfg.RetryCount)
	if err != nil {
		return nil, err
	}

	jsonResult := gjson.ParseBytes(resp.Body())
	out := map[string]interface{}{
		"price":                  jsonResult.Get("data.order.total_money").Str,
		"freight_discount_money": jsonResult.Get("data.order.freight_discount_money").Str,                // 运费折扣费用
		"freight_money":          jsonResult.Get("data.order.freight_money").Str,                         // 运费
		"order_freight":          jsonResult.Get("data.order.freights.0.freight.freight_real_money").Str, // 订单运费
		"user_ticket_id":         jsonResult.Get("data.order.default_coupon._id").Str,
	}
	return out, nil
}

func (s *Session) CreateOrder(
	ctx context.Context,
	cartData map[string]interface{},
	checkOrderData map[string]interface{},
	reserveTime *ReserveTime,
) error {
	urlPath := "https://maicai.api.ddxq.mobi/order/addNewOrder"

	paymentOrder := map[string]interface{}{
		"reserved_time_start":    reserveTime.StartTimestamp,
		"reserved_time_end":      reserveTime.EndTimestamp,
		"parent_order_sign":      cartData["parent_order_sign"],
		"address_id":             s.Address.Id,
		"pay_type":               s.PayType,
		"product_type":           1,
		"form_id":                strings.ReplaceAll(uuid.New().String(), "-", ""),
		"receipt_without_sku":    nil,
		"vip_money":              "",
		"vip_buy_user_ticket_id": "",
		"coupons_money":          "",
		"coupons_id":             "",
	}
	for k, v := range checkOrderData {
		paymentOrder[k] = v
	}

	packages := map[string]interface{}{
		"reserved_time_start":       reserveTime.StartTimestamp,
		"reserved_time_end":         reserveTime.EndTimestamp,
		"products":                  cartData["products"],
		"package_type":              cartData["package_type"],
		"package_id":                cartData["package_id"],
		"total_money":               cartData["total_money"],
		"total_origin_money":        cartData["total_origin_money"],
		"goods_real_money":          cartData["goods_real_money"],
		"total_count":               cartData["total_count"],
		"cart_count":                cartData["cart_count"],
		"is_presale":                cartData["is_presale"],
		"instant_rebate_money":      cartData["instant_rebate_money"],
		"coupon_rebate_money":       cartData["coupon_rebate_money"],
		"total_rebate_money":        cartData["total_rebate_money"],
		"used_balance_money":        cartData["used_balance_money"],
		"can_used_balance_money":    cartData["can_used_balance_money"],
		"used_point_num":            cartData["used_point_num"],
		"used_point_money":          cartData["used_point_money"],
		"can_used_point_num":        cartData["can_used_point_num"],
		"can_used_point_money":      cartData["can_used_point_money"],
		"is_share_station":          cartData["is_share_station"],
		"only_today_products":       cartData["only_today_products"],
		"only_tomorrow_products":    cartData["only_tomorrow_products"],
		"front_package_text":        cartData["front_package_text"],
		"front_package_type":        cartData["front_package_type"],
		"front_package_stock_color": cartData["front_package_stock_color"],
		"front_package_bg_color":    cartData["front_package_bg_color"],
		"eta_trace_id":              "",
		"soon_arrival":              "",
		"first_selected_big_time":   0,
		"receipt_without_sku":       0,
	}

	packageOrder := map[string]interface{}{
		"payment_order": paymentOrder,
		"packages":      []interface{}{packages},
	}
	packageOrderJson, err := json.Marshal(packageOrder)
	if err != nil {
		return fmt.Errorf("marshal products info failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Add("package_order", string(packageOrderJson))
	params.Add("showData", "true")
	params.Add("showMsg", "false")
	params.Add("ab_config", `{"key_onion":"C"}`)

	req := s.client.R()
	req.Header = s.buildHeader()
	req.SetBody(strings.NewReader(params.Encode()))
	_, err = s.execute(ctx, req, http.MethodPost, urlPath, s.cfg.RetryCount)
	return err
}
