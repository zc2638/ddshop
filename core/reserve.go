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
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

type ReserveTime struct {
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
	SelectMsg      string `json:"select_msg"`
}

func (s *Session) GetMultiReserveTime() ([]ReserveTime, error) {
	urlPath := "https://maicai.api.ddxq.mobi/order/getMultiReserveTime"
	var products []map[string]interface{}
	for _, product := range s.Order.Products {
		prod := map[string]interface{}{
			"id":                   product.Id,
			"total_money":          product.TotalPrice,
			"total_origin_money":   product.OriginPrice,
			"count":                product.Count,
			"price":                product.Price,
			"instant_rebate_money": "0.00",
			"origin_price":         product.OriginPrice,
		}
		products = append(products, prod)
	}
	productsList := [][]map[string]interface{}{
		products,
	}
	productsJson, err := json.Marshal(productsList)
	if err != nil {
		return nil, fmt.Errorf("marshal products info failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Add("group_config_id", "")
	params.Add("products", string(productsJson))
	params.Add("isBridge", "false")

	req := s.client.R()
	req.Header = s.buildHeader()
	req.SetBody(strings.NewReader(params.Encode()))
	resp, err := s.execute(context.Background(), req, http.MethodPost, urlPath)
	if err != nil {
		return nil, err
	}

	reserveTimes := gjson.Get(resp.String(), "data.0.time.0.times").Array()
	reserveTimeList := make([]ReserveTime, 0, len(reserveTimes))
	for _, reserveTimeInfo := range reserveTimes {
		if reserveTimeInfo.Get("disableType").Num != 0 {
			continue
		}
		reserveTime := ReserveTime{
			StartTimestamp: int(reserveTimeInfo.Get("start_timestamp").Num),
			EndTimestamp:   int(reserveTimeInfo.Get("end_timestamp").Num),
			SelectMsg:      reserveTimeInfo.Get("select_msg").Str,
		}
		reserveTimeList = append(reserveTimeList, reserveTime)
	}
	return reserveTimeList, nil
}
