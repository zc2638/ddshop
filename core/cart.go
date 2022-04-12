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
	"net/url"

	"github.com/tidwall/gjson"
)

// {"success":true,"code":0,"msg":"success","data":{"product":{"effective":[{"activity_info":{"id":"","gifts":null},"products":[{"id":"5e3f82cf7cdbf0131769408b","type":0,"category":"58fbf4fb936edf42508b4654","price":"4.59","sizes":[],"count":1,"status":1,"gifts":[],"addTime":1649606883,"cart_id":"5e3f82cf7cdbf0131769408b","activity_id":"","sku_activity_id":"","conditions_num":"","activity_tag":"","category_path":"58f9d213936edfe4568b569a,58fbf4fb936edf42508b4654","manage_category_path":"21,25,27","total_price":"4.59","origin_price":"4.59","no_supplementary_price":"4.59","no_supplementary_total_price":"4.59","size_price":"0.00","add_price":"4.59","add_vip_price":"","price_type":0,"buy_limit":0,"promotion_num":0,"product_name":"生姜 约300g","product_type":0,"small_image":"https://img.ddimg.mobi/product/3e7b7be5aa0b91616204086733.jpg?width=800&height=800","all_sizes":[],"only_new_user":false,"is_check":1,"is_gift":0,"is_bulk":0,"view_total_weight":"份","net_weight":"300","net_weight_unit":"g","is_stock":false,"old_count":1,"stock_number":1,"not_meet":[],"is_presale":0,"presale_id":"","presale_type":0,"delivery_start_time":0,"delivery_end_time":0,"is_invoice":1,"is_onion":0,"sub_list":[],"is_booking":0,"today_stockout":"","storage_value_id":0,"temperature_layer":"","is_shared_station_product":0,"is_fresh_food":0,"accessory_gifts":[],"accessory_text":"","supplementary_list":[]},{"id":"5e721d22b0055a0b5f763edf","type":0,"category":"58fbf4fb936edf42508b4654","price":"4.99","sizes":[],"count":1,"status":1,"gifts":[],"addTime":1649606846,"cart_id":"5e721d22b0055a0b5f763edf","activity_id":"","sku_activity_id":"","conditions_num":"","activity_tag":"","category_path":"58f9d213936edfe4568b569a,58fbf4fb936edf42508b4654","manage_category_path":"21,25,28","total_price":"4.99","origin_price":"4.99","no_supplementary_price":"4.99","no_supplementary_total_price":"4.99","size_price":"0.00","add_price":"4.99","add_vip_price":"","price_type":0,"buy_limit":0,"promotion_num":0,"product_name":"蒜头 约250g","product_type":0,"small_image":"https://img.ddimg.mobi/product/da62352cab2281613723470985.jpg?width=800&height=800","all_sizes":[],"only_new_user":false,"is_check":1,"is_gift":0,"is_bulk":0,"view_total_weight":"份","net_weight":"250","net_weight_unit":"g","is_stock":false,"old_count":1,"stock_number":1,"not_meet":[],"is_presale":0,"presale_id":"","presale_type":0,"delivery_start_time":0,"delivery_end_time":0,"is_invoice":1,"is_onion":0,"sub_list":[],"is_booking":0,"today_stockout":"","storage_value_id":0,"temperature_layer":"","is_shared_station_product":0,"is_fresh_food":0,"accessory_gifts":[],"accessory_text":"","supplementary_list":[]}]}],"invalid":[{"products":[{"id":"614d6cce8f1ed4f0871a2ca9","type":0,"category":"","price":"29.90","sizes":[],"count":1,"status":1,"gifts":[],"addTime":1649607493,"cart_id":"614d6cce8f1ed4f0871a2ca9","activity_id":"","sku_activity_id":"","conditions_num":"","activity_tag":"","category_path":"","manage_category_path":"258,259,262","origin_price":"29.90","size_price":"0.00","add_price":"29.90","add_vip_price":"","price_type":0,"buy_limit":0,"promotion_num":0,"product_name":"必品阁白菜猪肉王水饺 600g/袋","product_type":0,"small_image":"https://imgnew.ddimg.mobi/product/7f2617ebacf147999a4d356d375e6acf.gif?width=800&height=800","only_new_user":false,"is_check":0,"is_gift":0,"is_bulk":0,"view_total_weight":"袋","net_weight":"600","net_weight_unit":"g","is_stock":true,"old_count":1,"stock_number":0,"not_meet":[],"is_presale":0,"presale_id":"","presale_type":0,"delivery_start_time":0,"delivery_end_time":0,"is_invoice":1,"is_onion":0,"sub_list":[],"is_booking":0,"today_stockout":"","promotion_info":"","storage_value_id":3,"temperature_layer":"-18℃以下","is_fresh_food":0},{"id":"58ba8c02916edf9e4cc23072","type":0,"category":"58fb3b89936edfe4568b58ec","price":"9.90","sizes":[],"count":1,"status":1,"gifts":[],"addTime":1649607194,"cart_id":"58ba8c02916edf9e4cc23072","activity_id":"","sku_activity_id":"","conditions_num":"","activity_tag":"","category_path":"58f9e5a1936edf89778b568b,58fb3b89936edfe4568b58ec","manage_category_path":"330,331,332","origin_price":"9.90","size_price":"0.00","add_price":"9.90","add_vip_price":"","price_type":0,"buy_limit":0,"promotion_num":0,"product_name":"海天金标生抽酱油 500ml/瓶","product_type":0,"small_image":"https://ddimg.ddxq.mobi/879853186f70b1521771055327.jpg!maicai.product.list","only_new_user":false,"is_check":0,"is_gift":0,"is_bulk":0,"view_total_weight":"瓶","net_weight":"500","net_weight_unit":"ml","is_stock":true,"old_count":1,"stock_number":0,"not_meet":[],"is_presale":0,"presale_id":"","presale_type":0,"delivery_start_time":0,"delivery_end_time":0,"is_invoice":1,"is_onion":0,"sub_list":[],"is_booking":0,"today_stockout":"","promotion_info":"","storage_value_id":0,"temperature_layer":"","is_fresh_food":0}]}]},"toast":"","alert":null,"all_activity_cart":[],"station_id":"5c04bdd0716de1403a8b679b","order_product_list":[],"new_order_product_list":[{"products":[{"type":1,"id":"5e3f82cf7cdbf0131769408b","price":"4.59","count":1,"description":"","sizes":[],"cart_id":"5e3f82cf7cdbf0131769408b","parent_id":"","parent_batch_type":-1,"category_path":"58f9d213936edfe4568b569a,58fbf4fb936edf42508b4654","manage_category_path":"21,25,27","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"生姜 约300g","product_type":0,"small_image":"https://img.ddimg.mobi/product/3e7b7be5aa0b91616204086733.jpg?width=800&height=800","total_price":"4.59","origin_price":"4.59","total_origin_price":"4.59","no_supplementary_price":"4.59","no_supplementary_total_price":"4.59","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"份","net_weight":"300","net_weight_unit":"g","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":3,"is_presale":0},{"type":1,"id":"5e721d22b0055a0b5f763edf","price":"4.99","count":1,"description":"","sizes":[],"cart_id":"5e721d22b0055a0b5f763edf","parent_id":"","parent_batch_type":-1,"category_path":"58f9d213936edfe4568b569a,58fbf4fb936edf42508b4654","manage_category_path":"21,25,28","activity_id":"","sku_activity_id":"","conditions_num":"","product_name":"蒜头 约250g","product_type":0,"small_image":"https://img.ddimg.mobi/product/da62352cab2281613723470985.jpg?width=800&height=800","total_price":"4.99","origin_price":"4.99","total_origin_price":"4.99","no_supplementary_price":"4.99","no_supplementary_total_price":"4.99","size_price":"0.00","buy_limit":0,"price_type":0,"promotion_num":0,"instant_rebate_money":"0.00","is_invoice":1,"sub_list":[],"is_booking":0,"is_bulk":0,"view_total_weight":"份","net_weight":"250","net_weight_unit":"g","storage_value_id":0,"temperature_layer":"","sale_batches":{"batch_type":-1},"is_shared_station_product":0,"is_gift":0,"supplementary_list":[],"order_sort":4,"is_presale":0}],"total_money":"9.58","total_origin_money":"9.58","goods_real_money":"9.58","total_count":2,"cart_count":2,"is_presale":0,"instant_rebate_money":"0.00","used_balance_money":"0.00","can_used_balance_money":"0.00","used_point_num":0,"used_point_money":"0.00","can_used_point_num":0,"can_used_point_money":"0.00","is_share_station":0,"only_today_products":[],"only_tomorrow_products":[],"package_type":1,"package_id":1,"front_package_text":"即时配送","front_package_type":0,"front_package_stock_color":"#2FB157","front_package_bg_color":"#fbfefc"}],"order_product_list_sign":"d751713988987e9331980363e24189ce","full_to_off":"0.00","freight_money":"0.00","free_freight_type":3,"instant_rebate_money":"0.00","goods_real_money":"9.58","total_money":"9.58","is_select_detail":1,"good_max_count_toast":"订单商品明细行数超过最大限制，无法按商品明细开票","is_all_check":1,"onion_id":"","onion_tip":{"tip_name_type":0,"tip_name":"赠品小葱已赠完，如有需要可购买小葱","event_track_type":9},"cart_notice":"已免配送费","cart_notice_new":"免配送费","free_freight_notice":{},"cart_top_floor_info":[],"cart_count":2,"total_count":4,"product_num":{"5e721d22b0055a0b5f763edf":1,"614d6cce8f1ed4f0871a2ca9":1,"5e3f82cf7cdbf0131769408b":1,"58ba8c02916edf9e4cc23072":1},"stop_order_toast":"","gift_no_size_tip":"","is_hit_onion":false,"onion_ab_config":3,"is_hit_gift_size":true,"coupon_text_a":"","coupon_text_b":"","need_amount":"","is_vip_ticket":0,"coupon_amount":"","coupon_state":-1,"coupon_type":0,"next_recommend_coupon":{"coupon_text_a":null,"coupon_text_b":null,"need_amount":null,"is_vip_ticket":null,"is_common_ticket":null},"show_coupon_detail":false,"contains_advent_gift":0,"parent_order_info":{"parent_order_sign":"5192235f19162dbe7f1aa1cf749717ba","stockout_gift_product":[],"stockout_gift_text":"赠品赠完即止，不再补送，敬请谅解。","is_open_presale_use_virtual_stock":false},"is_support_merge_payment":1,"sodexo_nonsupport_product_list":[],"valid_product_counts":{"5e721d22b0055a0b5f763edf":1,"5e3f82cf7cdbf0131769408b":1}},"tradeTag":"success","server_time":1649627313,"is_trade":1}

type ProductResult struct {
	Error      string      `json:"error"`
	Success    bool        `json:"success"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       ProductData `json:"data"`
	TradeTag   string      `json:"tradeTag"`
	ServerTime int         `json:"server_time"`
	IsTrade    int         `json:"is_trade"`
}

type ProductData struct {
	Product              ProductSet         `json:"product"`
	Toast                string             `json:"toast"`
	Alert                interface{}        `json:"alert"`
	AllActivityCart      []interface{}      `json:"all_activity_cart"`
	StationId            string             `json:"station_id"`
	OrderProductList     []interface{}      `json:"order_product_list"`
	NewOrderProductList  []OrderProductInfo `json:"new_order_product_list"`
	OrderProductListSign string             `json:"order_product_list_sign"`
	FullToOff            string             `json:"full_to_off"`
	FreightMoney         string             `json:"freight_money"`
	FreeFreightType      int                `json:"free_freight_type"`
	InstantRebateMoney   string             `json:"instant_rebate_money"`
	GoodsRealMoney       string             `json:"goods_real_money"`
	TotalMoney           string             `json:"total_money"`
	IsSelectDetail       int                `json:"is_select_detail"`
	GoodMaxCountToast    string             `json:"good_max_count_toast"`
	IsAllCheck           int                `json:"is_all_check"`
	OnionId              string             `json:"onion_id"`
	OnionTip             struct{}           `json:"onion_tip"`
	CartNotice           string             `json:"cart_notice"`
	CartNoticeNew        string             `json:"cart_notice_new"`
	FreeFreightNotice    struct{}           `json:"free_freight_notice"`
	CartTopFloorInfo     []interface{}      `json:"cart_top_floor_info"`
	CartCount            int                `json:"cart_count"`
	TotalCount           int                `json:"total_count"`
	StopOrderToast       string             `json:"stop_order_toast"`
}

type ProductSet struct {
	Effective []ProductInfo `json:"effective"`
	Invalid   []ProductInfo `json:"invalid"`
}

type ProductInfo struct {
	Products []Product `json:"products"`
}

type Product struct {
	Id                        string                   `json:"id"`
	Type                      int                      `json:"type"`
	Category                  string                   `json:"category"`
	Price                     string                   `json:"price"`
	Sizes                     []map[string]interface{} `json:"sizes"`
	Count                     int                      `json:"count"`
	Status                    int                      `json:"status"`
	Gifts                     []interface{}            `json:"gifts"`
	AddTime                   int                      `json:"addTime"`
	CartId                    string                   `json:"cart_id"`
	ActivityId                string                   `json:"activity_id"`
	SkuActivityId             string                   `json:"sku_activity_id"`
	ConditionsNum             string                   `json:"conditions_num"`
	ActivityTag               string                   `json:"activity_tag"`
	CategoryPath              string                   `json:"category_path"`
	ManageCategoryPath        string                   `json:"manage_category_path"`
	TotalPrice                string                   `json:"total_price"`
	OriginPrice               string                   `json:"origin_price"`
	NoSupplementaryPrice      string                   `json:"no_supplementary_price"`
	NoSupplementaryTotalPrice string                   `json:"no_supplementary_total_price"`
	SizePrice                 string                   `json:"size_price"`
	AddPrice                  string                   `json:"add_price"`
	AddVipPrice               string                   `json:"add_vip_price"`
	PriceType                 int                      `json:"price_type"`
	BuyLimit                  int                      `json:"buy_limit"`
	PromotionNum              int                      `json:"promotion_num"`
	ProductName               string                   `json:"product_name"`
	ProductType               int                      `json:"product_type"`
	SmallImage                string                   `json:"small_image"`
	AllSizes                  []interface{}            `json:"all_sizes"`
	OnlyNewUser               bool                     `json:"only_new_user"`
	IsCheck                   int                      `json:"is_check"`
	IsGift                    int                      `json:"is_gift"`
	IsBulk                    int                      `json:"is_bulk"`
	ViewTotalWeight           string                   `json:"view_total_weight"`
	NetWeight                 string                   `json:"net_weight"`
	NetWeightUnit             string                   `json:"net_weight_unit"`
	IsStock                   bool                     `json:"is_stock"`
	OldCount                  int                      `json:"old_count"`
	StockNumber               int                      `json:"stock_number"`
	NotMeet                   []interface{}            `json:"not_meet"`
	IsPresale                 int                      `json:"is_presale"`
	PresaleId                 string                   `json:"presale_id"`
	PresaleType               int                      `json:"presale_type"`
	DeliveryStartTime         int                      `json:"delivery_start_time"`
	DeliveryEndTime           int                      `json:"delivery_end_time"`
	IsInvoice                 int                      `json:"is_invoice"`
	IsOnion                   int                      `json:"is_onion"`
	SubList                   []interface{}            `json:"sub_list"`
	IsBooking                 int                      `json:"is_booking"`
	TodayStockout             string                   `json:"today_stockout"`
	StorageValueId            int                      `json:"storage_value_id"`
	TemperatureLayer          string                   `json:"temperature_layer"`
	IsSharedStationProduct    int                      `json:"is_shared_station_product"`
	IsFreshFood               int                      `json:"is_fresh_food"`
	AccessoryGifts            []interface{}            `json:"accessory_gifts"`
	AccessoryText             string                   `json:"accessory_text"`
	SupplementaryList         []interface{}            `json:"supplementary_list"`
}

type OrderProductInfo struct {
	Products               []Product     `json:"products"`
	TotalMoney             string        `json:"total_money"`
	TotalOriginMoney       string        `json:"total_origin_money"`
	GoodsRealMoney         string        `json:"goods_real_money"`
	TotalCount             int           `json:"total_count"`
	CartCount              int           `json:"cart_count"`
	IsPresale              int           `json:"is_presale"`
	InstantRebateMoney     string        `json:"instant_rebate_money"`
	TotalRebateMoney       string        `json:"total_rebate_money"`
	UsedBalanceMoney       string        `json:"used_balance_money"`
	CanUsedBalanceMoney    string        `json:"can_used_balance_money"`
	UsedPointNum           int           `json:"used_point_num"`
	UsedPointMoney         string        `json:"used_point_money"`
	CanUsedPointNum        int           `json:"can_used_point_num"`
	CanUsedPointMoney      string        `json:"can_used_point_money"`
	IsShareStation         int           `json:"is_share_station"`
	OnlyTodayProducts      []interface{} `json:"only_today_products"`
	OnlyTomorrowProducts   []interface{} `json:"only_tomorrow_products"`
	PackageType            int           `json:"package_type"`
	PackageId              int           `json:"package_id"`
	FrontPackageText       string        `json:"front_package_text"`
	FrontPackageType       int           `json:"front_package_type"`
	FrontPackageStockColor string        `json:"front_package_stock_color"`
	FrontPackageBgColor    string        `json:"front_package_bg_color"`
}

type OrderProduct struct {
	Type                      int                      `json:"type"`
	Id                        string                   `json:"id"`
	Price                     string                   `json:"price"`
	Count                     int                      `json:"count"`
	Description               string                   `json:"description"`
	Sizes                     []map[string]interface{} `json:"sizes"`
	CartId                    string                   `json:"cart_id"`
	ParentId                  string                   `json:"parent_id"`
	ParentBatchType           int                      `json:"parent_batch_type"`
	CategoryPath              string                   `json:"category_path"`
	ManageCategoryPath        string                   `json:"manage_category_path"`
	ActivityId                string                   `json:"activity_id"`
	SkuActivityId             string                   `json:"sku_activity_id"`
	ConditionsNum             string                   `json:"conditions_num"`
	ProductName               string                   `json:"product_name"`
	ProductType               int                      `json:"product_type"`
	SmallImage                string                   `json:"small_image"`
	TotalPrice                string                   `json:"total_price"`
	OriginPrice               string                   `json:"origin_price"`
	TotalOriginPrice          string                   `json:"total_origin_price"`
	NoSupplementaryPrice      string                   `json:"no_supplementary_price"`
	NoSupplementaryTotalPrice string                   `json:"no_supplementary_total_price"`
	SizePrice                 string                   `json:"size_price"`
	BuyLimit                  int                      `json:"buy_limit"`
	PriceType                 int                      `json:"price_type"`
	PromotionNum              int                      `json:"promotion_num"`
	InstantRebateMoney        string                   `json:"instant_rebate_money"`
	IsInvoice                 int                      `json:"is_invoice"`
	SubList                   []interface{}            `json:"sub_list"`
	IsBooking                 int                      `json:"is_booking"`
	IsBulk                    int                      `json:"is_bulk"`
	ViewTotalWeight           string                   `json:"view_total_weight"`
	NetWeight                 string                   `json:"net_weight"`
	NetWeightUnit             string                   `json:"net_weight_unit"`
	StorageValueId            int                      `json:"storage_value_id"`
	TemperatureLayer          string                   `json:"temperature_layer"`
	SaleBatches               struct {
		BatchType int `json:"batch_type"`
	} `json:"sale_batches"`
	IsSharedStationProduct int           `json:"is_shared_station_product"`
	IsGift                 int           `json:"is_gift"`
	SupplementaryList      []interface{} `json:"supplementary_list"`
	OrderSort              int           `json:"order_sort"`
	IsPresale              int           `json:"is_presale"`
}

type Cart struct {
	ProdList        []Product `json:"effective_products"`
	ParentOrderSign string    `json:"parent_order_sign"`
}

func (s *Session) CartAllCheck() error {
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
	_, err = s.execute(context.Background(), req, http.MethodGet, urlPath)
	return err
}

func (s *Session) GetCart() error {
	u, err := url.Parse("https://maicai.api.ddxq.mobi/cart/index")
	if err != nil {
		return fmt.Errorf("cart url parse failed: %v", err)
	}

	params := s.buildURLParams(true)
	params.Set("is_load", "1")
	params.Set("ab_config", "{\"key_onion\":\"D\",\"key_cart_discount_price\":\"C\"}")
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	req.Header = s.buildHeader()
	resp, err := s.execute(context.Background(), req, http.MethodGet, urlPath)
	if err != nil {
		return err
	}

	var productResult ProductResult
	if err := json.Unmarshal(resp.Body(), &productResult); err != nil {
		return fmt.Errorf("parse response failed: %v, body: %v", err, resp.String())
	}

	jsonResult := gjson.ParseBytes(resp.Body())
	s.Cart.ParentOrderSign = jsonResult.Get("data.parent_order_info.parent_order_sign").Str
	switch s.CartMode {
	case 1:
		var products []Product
		for _, v := range productResult.Data.Product.Effective {
			products = append(products, v.Products...)
		}
		s.Cart.ProdList = products
	case 2:
		var products []Product
		for _, v := range productResult.Data.NewOrderProductList {
			products = append(products, v.Products...)
		}
		s.Cart.ProdList = products
	default:
		return fmt.Errorf("incorrect cart mode: %v", s.CartMode)
	}
	return nil
}
