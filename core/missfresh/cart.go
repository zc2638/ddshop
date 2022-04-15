package missfresh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CartResult struct {
	Data    CartData `json:"data"`
	Code    int      `json:"code"`
	Msg     string   `json:"msg"`
	Success bool     `json:"success"`
}

type CartData struct {
	AllChecked int `json:"allChecked"`
	AbTest     struct {
		CartBatchSimilar int `json:"cartBatchSimilar"`
	} `json:"abTest"`
	ShowResourceSpace int `json:"showResourceSpace"`
	BalanceArea       struct {
		BalanceCount int    `json:"balanceCount"`
		PayTips      string `json:"payTips"`
		PayAmount    int    `json:"payAmount"`
		CartSkuCount int    `json:"cartSkuCount"`
	} `json:"balanceArea"`
	SaleGroups []struct {
		SellerType          int                    `json:"sellerType"`
		SaleGroupId         string                 `json:"saleGroupId"`
		SaleGroupType       string                 `json:"saleGroupType"`
		SaleGroupTitleImage string                 `json:"saleGroupTitleImage"`
		SaleGroupTitle      string                 `json:"saleGroupTitle"`
		SaleGroupChecked    int                    `json:"saleGroupChecked"`
		ExchangeBuy         map[string]interface{} `json:"exchangeBuy"`
		SaleGroupPostage    map[string]interface{} `json:"saleGroupPostage"`
		PromGroups          [][]struct {
			ShowStyle     int `json:"showStyle"`
			PromGroupHead struct {
				PromGroupId    string `json:"promGroupId"`
				PromGroupTag   string `json:"promGroupTag"`
				PromGroupTitle string `json:"promGroupTitle"`
				ToAddOnText    string `json:"toAddOnText"`
				ToAddOnUrl     string `json:"toAddOnUrl"`
				MakeOrder      struct {
					MakeOrderType     int    `json:"makeOrderType"`
					PromotionId       int    `json:"promotionId"`
					UnifiedCategoryId int    `json:"unifiedCategoryId"`
					SaleGroupTypeStr  string `json:"saleGroupTypeStr"`
				} `json:"makeOrder"`
				BuryParam struct {
					StrategyId string `json:"strategy_id"`
				} `json:"buryParam"`
				BuryItem struct {
					ToAddOnType string `json:"toAddOnType"`
				} `json:"buryItem"`
			} `json:"promGroupHead,omitempty"`
			Product struct {
				Sku           string        `json:"sku"`
				Name          string        `json:"name"`
				Image         string        `json:"image"`
				SubTitles     []interface{} `json:"subTitles"`
				PromotionTips string        `json:"promotionTips"`
				Quantity      int           `json:"quantity"`
				Stock         int           `json:"stock"`
				Price         struct {
					TotalSellPrice struct {
						ShowStyle int `json:"showStyle"`
						Color     int `json:"color"`
						Price     int `json:"price"`
					} `json:"totalSellPrice"`
					SellPrice struct {
						ShowStyle int `json:"showStyle"`
						Color     int `json:"color"`
						Price     int `json:"price"`
					} `json:"sellPrice"`
				} `json:"price"`
				Type            string `json:"type"`
				Unit            string `json:"unit"`
				Status          int    `json:"status"`
				SaleGroupId     string `json:"saleGroupId"`
				PromGroupId     string `json:"promGroupId"`
				Checked         int    `json:"checked"`
				ShowChecked     int    `json:"showChecked"`
				SupportedOption int    `json:"supportedOption"`
				DeliveryTag     []struct {
					TagName     string `json:"tagName"`
					TextColor   int    `json:"textColor"`
					BgColor     int    `json:"bgColor"`
					BorderColor int    `json:"borderColor"`
					UseType     int    `json:"useType"`
				} `json:"deliveryTag"`
				PromotionDetails []interface{} `json:"promotionDetails"`
				PurchaseLimit    int           `json:"purchaseLimit"`
				BuryParam        struct {
					SkuStore    string `json:"skuStore"`
					SkuSupplier string `json:"skuSupplier"`
				} `json:"buryParam"`
			} `json:"product,omitempty"`
		} `json:"promGroups"`
		VoucherList []interface{} `json:"voucherList"`
	} `json:"saleGroups"`
	PostageRules []struct {
		PostageRuleId string `json:"postageRuleId"`
		Title         string `json:"title"`
		Rules         string `json:"rules"`
	} `json:"postageRules"`
}

func (s *Session) GetCart(ctx context.Context) error {
	u, err := url.Parse(HOST + CART_PROMOTION_API)
	if err != nil {
		return fmt.Errorf("cart url parse failed: %v", err)
	}

	urlPath := u.String()

	req := s.client.R()
	resp, err := s.execute(ctx, req, http.MethodPost, urlPath, maxRetryCount)
	if err != nil {
		return err
	}

	var cartResult CartResult
	if err := json.Unmarshal(resp.Body(), &cartResult); err != nil {
		return fmt.Errorf("parse response failed: %v", err)
	}

	return nil
}
