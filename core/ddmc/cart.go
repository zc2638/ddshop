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
	"net/url"
)

func (s *Session) CartAllCheck(ctx context.Context) error {
	u, err := url.Parse("https://maicai.api.ddxq.mobi/cart/allCheck")
	if err != nil {
		return fmt.Errorf("cart url parse failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Set("is_check", "1")
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	req.Header = s.buildHeader()
	_, err = s.execute(ctx, req, http.MethodGet, urlPath, maxRetryCount)
	return err
}

func (s *Session) GetCart(ctx context.Context) (map[string]interface{}, error) {
	u, err := url.Parse("https://maicai.api.ddxq.mobi/cart/index")
	if err != nil {
		return nil, fmt.Errorf("获取购物车商品，请求URL解析失败: %v", err)
	}

	params := s.buildURLParams(true)
	params.Set("is_load", "1")
	params.Set("ab_config", `{"key_onion":"D","key_cart_discount_price":"C"}`)
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	req.Header = s.buildHeader()
	resp, err := s.execute(ctx, req, http.MethodGet, urlPath, maxRetryCount)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("获取购物车商品，JSON解析失败: %s", resp.String())
	}

	data := result["data"].(map[string]interface{})
	list := data["new_order_product_list"].([]interface{})
	if len(list) == 0 {
		return nil, ErrorNoValidProduct
	}

	item := list[0].(map[string]interface{})

	productList := item["products"].([]interface{})
	products := make([]map[string]interface{}, 0, len(productList))
	for _, v := range productList {
		product := v.(map[string]interface{})
		product["total_money"] = product["total_price"]
		product["total_origin_money"] = product["total_origin_price"]
		products = append(products, product)
	}

	out := map[string]interface{}{
		"products":                  products,
		"package_type":              item["package_type"],
		"package_id":                item["package_id"],
		"total_money":               item["total_money"],
		"total_origin_money":        item["total_origin_money"],
		"goods_real_money":          item["goods_real_money"],
		"total_count":               item["total_count"],
		"cart_count":                item["cart_count"],
		"is_presale":                item["is_presale"],
		"instant_rebate_money":      item["instant_rebate_money"],
		"coupon_rebate_money":       item["coupon_rebate_money"],
		"total_rebate_money":        item["total_rebate_money"],
		"used_balance_money":        item["used_balance_money"],
		"can_used_balance_money":    item["can_used_balance_money"],
		"used_point_num":            item["used_point_num"],
		"used_point_money":          item["used_point_money"],
		"can_used_point_num":        item["can_used_point_num"],
		"can_used_point_money":      item["can_used_point_money"],
		"is_share_station":          item["is_share_station"],
		"only_today_products":       item["only_today_products"],
		"only_tomorrow_products":    item["only_tomorrow_products"],
		"front_package_text":        item["front_package_text"],
		"front_package_type":        item["front_package_type"],
		"front_package_stock_color": item["front_package_stock_color"],
		"front_package_bg_color":    item["front_package_bg_color"],
	}
	parentOrderInfo, ok := data["parent_order_info"].(map[string]interface{})
	if ok {
		out["parent_order_sign"] = parentOrderInfo["parent_order_sign"]
	}
	return out, nil
}
