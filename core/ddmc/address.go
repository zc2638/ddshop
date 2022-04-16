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
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type AddressResult struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    AddressData `json:"data"`
}

type Address struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	StationId  string  `json:"station_id"`
	CityNumber string  `json:"city_number"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	UserName   string  `json:"user_name"`
	Mobile     string  `json:"mobile"`
	Address    string  `json:"address"`
	AddrDetail string  `json:"addr_detail"`
}

type AddressData struct {
	ValidAddress    []AddressItem `json:"valid_address"`
	InvalidAddress  []AddressItem `json:"invalid_address"`
	MaxAddressCount int           `json:"max_address_count"`
	CanAddAddress   bool          `json:"can_add_address"`
}

type AddressItem struct {
	Id          string             `json:"id"`
	Gender      int                `json:"gender"`
	Mobile      string             `json:"mobile"`
	Location    AddressLocation    `json:"location"`
	Label       string             `json:"label"`
	UserName    string             `json:"user_name"`
	AddrDetail  string             `json:"addr_detail"`
	StationId   string             `json:"station_id"`
	StationName string             `json:"station_name"`
	IsDefault   bool               `json:"is_default"`
	CityNumber  string             `json:"city_number"`
	InfoStatus  int                `json:"info_status"`
	StationInfo AddressStationInfo `json:"station_info"`
	VillageId   string             `json:"village_id"`
}

type AddressLocation struct {
	TypeCode string    `json:"typecode"`
	Address  string    `json:"address"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
	Id       string    `json:"id"`
}

type AddressStationInfo struct {
	Id           string `json:"id"`
	Address      string `json:"address"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	BusinessTime string `json:"business_time"`
	CityName     string `json:"city_name"`
	CityNumber   string `json:"city_number"`
}

func (s *Session) GetAddress() (map[string]AddressItem, error) {
	u, err := url.Parse("https://sunquan.api.ddxq.mobi/api/v1/user/address/")
	if err != nil {
		return nil, fmt.Errorf("address url parse failed: %v", err)
	}

	params := s.buildURLParams(false)
	params.Set("source_type", "5")
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	req.SetHeader("Host", "sunquan.api.ddxq.mobi")
	resp, err := s.execute(context.Background(), req, http.MethodGet, urlPath, s.cfg.RetryCount)
	if err != nil {
		return nil, err
	}

	var addressResult AddressResult
	if err := json.Unmarshal(resp.Body(), &addressResult); err != nil {
		return nil, fmt.Errorf("parse response failed: %v", err)
	}
	if len(addressResult.Data.ValidAddress) == 0 {
		return nil, errors.New("未查询到有效收货地址，请前往 app 添加或检查填写的 cookie 是否正确！")
	}

	result := make(map[string]AddressItem)
	for _, v := range addressResult.Data.ValidAddress {
		str := fmt.Sprintf("%s %s %s", v.UserName, v.Location.Address, v.AddrDetail)
		result[str] = v
	}
	return result, nil
}
