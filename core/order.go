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

type Order struct {
	Products []Product `json:"products"`
	Price    string    `json:"price"`
}

type Package struct {
	FirstSelectedBigTime string                   `json:"first_selected_big_time"`
	Products             []map[string]interface{} `json:"products"`
	EtaTraceId           string                   `json:"eta_trace_id"`
	PackageId            int                      `json:"package_id"`
	ReservedTimeStart    int                      `json:"reserved_time_start"`
	ReservedTimeEnd      int                      `json:"reserved_time_end"`
	SoonArrival          int                      `json:"soon_arrival"`
	PackageType          int                      `json:"package_type"`
}

type PaymentOrder struct {
	ReservedTimeStart    int    `json:"reserved_time_start"`
	ReservedTimeEnd      int    `json:"reserved_time_end"`
	FreightDiscountMoney string `json:"freight_discount_money"`
	FreightMoney         string `json:"freight_money"`
	OrderFreight         string `json:"order_freight"`
	AddressId            string `json:"address_id"`
	UsedPointNum         int    `json:"used_point_num"`
	ParentOrderSign      string `json:"parent_order_sign"`
	PayType              int    `json:"pay_type"`
	OrderType            int    `json:"order_type"`
	IsUseBalance         int    `json:"is_use_balance"`
	ReceiptWithoutSku    string `json:"receipt_without_sku"`
	Price                string `json:"price"`
}

type PackageOrder struct {
	Packages     []*Package   `json:"packages"`
	PaymentOrder PaymentOrder `json:"payment_order"`
}

type AddNewOrderReturnData struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    struct {
		PackageOrder     PackageOrder `json:"package_order"`
		StockoutProducts []Product    `json:"stockout_products"`
	} `json:"data"`
}

func (s *Session) GeneratePackageOrder() {
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
			"sizes":                product.Sizes,
		}
		products = append(products, prod)
	}

	p := Package{
		FirstSelectedBigTime: "0",
		Products:             products,
		EtaTraceId:           "",
		PackageId:            1,
		PackageType:          1,
	}
	paymentOrder := PaymentOrder{
		FreightDiscountMoney: "5.00",
		FreightMoney:         "5.00",
		OrderFreight:         "0.00",
		AddressId:            s.Address.Id,
		UsedPointNum:         0,
		ParentOrderSign:      s.Cart.ParentOrderSign,
		PayType:              s.PayType,
		OrderType:            1,
		IsUseBalance:         0,
		ReceiptWithoutSku:    "1",
		Price:                s.Order.Price,
	}
	packageOrder := PackageOrder{
		Packages: []*Package{
			&p,
		},
		PaymentOrder: paymentOrder,
	}
	s.PackageOrder = packageOrder
}

func (s *Session) UpdatePackageOrder(reserveTime ReserveTime) {
	s.PackageOrder.PaymentOrder.ReservedTimeStart = reserveTime.StartTimestamp
	s.PackageOrder.PaymentOrder.ReservedTimeEnd = reserveTime.EndTimestamp
	for _, p := range s.PackageOrder.Packages {
		p.ReservedTimeStart = reserveTime.StartTimestamp
		p.ReservedTimeEnd = reserveTime.EndTimestamp
	}
}

func (s *Session) CheckOrder() error {
	urlPath := "https://maicai.api.ddxq.mobi/order/checkOrder"

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
			"sizes":                product.Sizes,
		}
		products = append(products, prod)
	}
	packagesInfo := []map[string]interface{}{
		{
			"package_type": 1,
			"package_id":   1,
			"products":     products,
		},
	}
	packagesJson, err := json.Marshal(packagesInfo)
	if err != nil {
		return fmt.Errorf("marshal products info failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Add("user_ticket_id", "default")
	params.Add("freight_ticket_id", "default")
	params.Add("is_use_point", "0")
	params.Add("is_use_balance", "0")
	params.Add("is_buy_vip", "0")
	params.Add("coupons_id", "")
	params.Add("is_buy_coupons", "0")
	params.Add("packages", string(packagesJson))
	params.Add("check_order_type", "0")
	params.Add("is_support_merge_payment", "1")
	params.Add("showData", "true")
	params.Add("showMsg", "false")

	req := s.client.R()
	req.Header = s.buildHeader()
	req.SetBody(strings.NewReader(params.Encode()))
	resp, err := s.execute(context.TODO(), req, http.MethodPost, urlPath)
	if err != nil {
		return err
	}

	s.Order.Price = gjson.Get(resp.String(), "data.order.total_money").Str
	return nil
}

func (s *Session) CreateOrder(ctx context.Context) error {
	urlPath := "https://maicai.api.ddxq.mobi/order/addNewOrder"

	packageOrderJson, err := json.Marshal(s.PackageOrder)
	if err != nil {
		return fmt.Errorf("marshal products info failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Add("package_order", string(packageOrderJson))
	params.Add("showData", "true")
	params.Add("showMsg", "false")
	params.Add("ab_config", "{\"key_onion\":\"C\"}")

	req := s.client.R()
	req.Header = s.buildHeader()
	req.SetBody(strings.NewReader(params.Encode()))
	_, err = s.execute(ctx, req, http.MethodPost, urlPath)
	return err
}
