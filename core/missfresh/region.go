package missfresh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type RegionResult struct {
	Data    RegionData `json:"data"`
	Code    int        `json:"code"`
	Msg     string     `json:"msg"`
	Success bool       `json:"success"`
}

type RegionData struct {
	Lat                 string `json:"lat"`
	Lng                 string `json:"lng"`
	AreaCode            string `json:"areaCode"`
	City                string `json:"city"`
	StationName         string `json:"stationName"`
	BigWarehouse        string `json:"bigWarehouse"`
	WarehouseCode       string `json:"warehouseCode"`
	StationCode         string `json:"stationCode"`
	ImgUrl              string `json:"imgUrl"`
	WhiteChromeImageUrl string `json:"whiteChromeImageUrl"`
	AssistLocationFlag  int    `json:"assistLocationFlag"`
	RegionId            int    `json:"regionId"`
	Type                int    `json:"type"`
	DeliveryType        int    `json:"deliveryType"`
	SellerId            int    `json:"sellerId"`
	SellerInfoList      []struct {
		SellerId   int `json:"sellerId"`
		SellerType int `json:"sellerType"`
	} `json:"sellerInfoList"`
	Tips string `json:"tips"`
}

type RegionInfo struct {
	Lat           string `json:"lat"`
	Lng           string `json:"lng"`
	AreaCode      string `json:"areaCode"`
	BigWarehouse  string `json:"bigWarehouse"`
	WarehouseCode string `json:"warehouseCode"`
	StationCode   string `json:"stationCode"`
	RegionId      int    `json:"regionId"`
	Type          int    `json:"type"`
	DeliveryType  int    `json:"deliveryType"`
	SellerId      int    `json:"sellerId"`
}

func (s *Session) GetRegion(ctx context.Context) error {
	u, err := url.Parse(HOST + REGION_API)
	if err != nil {
		return fmt.Errorf("region url parse failed: %v", err)
	}

	urlPath := u.String()

	req := s.client.R()
	body := s.CreateBodyParam()
	i := body["common"].(map[string]interface{})
	i["addressCode"] = s.DefaultAddressInfo.AreaCode
	i["userLat"] = s.DefaultAddressInfo.Lat
	i["userLng"] = s.DefaultAddressInfo.Lng
	i["currentLat"] = s.DefaultAddressInfo.Lat
	i["currentLng"] = s.DefaultAddressInfo.Lng
	body["param"] = map[string]interface{}{
		"lat":       s.DefaultAddressInfo.Lat,
		"lng":       s.DefaultAddressInfo.Lng,
		"addressId": s.DefaultAddressInfo.AddressID,
	}
	body["common"] = i
	marshal, err := json.Marshal(body)
	if err != nil {
		return err
	}
	fmt.Println(string(marshal))
	req.SetBody(marshal)
	req.SetHeader("x-region", fmt.Sprintf(`{"address_code":"%s","station_code":""}`, s.DefaultAddressInfo.AreaCode))

	resp, err := s.execute(ctx, req, http.MethodPost, urlPath, maxRetryCount)
	if err != nil {
		return err
	}

	var regionResult RegionResult
	if err := json.Unmarshal(resp.Body(), &regionResult); err != nil {
		return fmt.Errorf("parse response failed: %v", err)
	}

	s.RegionInfo = &RegionInfo{
		Lat:           regionResult.Data.Lat,
		Lng:           regionResult.Data.Lng,
		AreaCode:      regionResult.Data.AreaCode,
		BigWarehouse:  regionResult.Data.BigWarehouse,
		WarehouseCode: regionResult.Data.WarehouseCode,
		StationCode:   regionResult.Data.StationCode,
		RegionId:      regionResult.Data.RegionId,
		Type:          regionResult.Data.Type,
		DeliveryType:  regionResult.Data.DeliveryType,
		SellerId:      regionResult.Data.SellerId,
	}
	return nil
}
