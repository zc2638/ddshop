package missfresh

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type AddressResult struct {
	Addresses []AddressItem `json:"addresses"`
	Code      int           `json:"code"`
}

type AddressItem struct {
	Gender        int    `json:"gender"`
	Optimal       int    `json:"optimal"`
	City          string `json:"city"`
	Address1      string `json:"address_1"`
	Address2      string `json:"address_2"`
	AreaCode      string `json:"area_code"`
	FullAddress   string `json:"full_address"`
	Transport     bool   `json:"transport"`
	AddressDetail string `json:"address_detail"`
	LatLng        string `json:"lat_lng"`
	Province      string `json:"province"`
	Name          string `json:"name"`
	PhoneNumber   string `json:"phone_number"`
	Id            int    `json:"id"`
	Tag           string `json:"tag"`
	PoiId         string `json:"poi_id"`
	Area          string `json:"area,omitempty"`
}

func (s *Session) GetAddress(ctx context.Context) (map[string]AddressItem, error) {
	u, err := url.Parse(HOST + ADDRESS_API)
	if err != nil {
		return nil, fmt.Errorf("address url parse failed: %v", err)
	}

	params := s.buildURLParams(false)
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	resp, err := s.execute(ctx, req, http.MethodGet, urlPath, maxRetryCount)
	if err != nil {
		return nil, err
	}

	var addressResult AddressResult
	if err := json.Unmarshal(resp.Body(), &addressResult); err != nil {
		return nil, fmt.Errorf("parse response failed: %v", err)
	}
	if len(addressResult.Addresses) == 0 {
		return nil, errors.New("未查询到有效收货地址，请前往 app 添加或检查填写的 cookie 是否正确！")
	}

	result := make(map[string]AddressItem)
	for _, v := range addressResult.Addresses {
		str := fmt.Sprintf("%s %s %s", v.Name, v.AreaCode, v.AddressDetail)
		result[str] = v
	}
	return result, nil
}
