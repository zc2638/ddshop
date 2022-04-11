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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

/**
{
    "success": true,
    "code": 0,
    "message": "",
    "data": {
        "valid_address": [
            {
                "id": "6252ae9f5847f50001389f94",
                "gender": 1,
                "mobile": "138****001",
                "location": {
                    "typecode": "120302",
                    "address": "宝山区殷高路7弄(殷高路地铁站3号口步行290米)",
                    "name": "殷高路7弄小区",
                    "location": [
                        121.493507,
                        31.321424
                    ],
                    "id": "B0FFHUBV50"
                },
                "label": "",
                "user_name": "郑",
                "addr_detail": "xxxxx",
                "station_id": "5c04bdd0716de1403a8b679b",
                "station_name": "高境站",
                "is_default": true,
                "city_number": "0101",
                "info_status": 1,
                "station_info": {
                    "id": "5c04bdd0716de1403a8b679b",
                    "address": "",
                    "name": "高境站",
                    "phone": "10103365",
                    "business_time": "24h",
                    "city_name": "上海市",
                    "city_number": "0101"
                },
                "village_id": "5ec781010ae80b6400a1c156"
            }
        ],
        "invalid_address": [],
        "max_address_count": 10,
        "can_add_address": true
    }
}
*/

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
	resp, err := req.Get(urlPath)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("statusCode: %d, body: %s", resp.StatusCode(), resp.String())
	}

	var addressResult AddressResult
	if err := json.Unmarshal(resp.Body(), &addressResult); err != nil {
		return nil, fmt.Errorf("parse response failed: %v", err)
	}
	if addressResult.Code != 0 || !addressResult.Success {
		return nil, fmt.Errorf("request address failed: %v", addressResult.Message)
	}
	if len(addressResult.Data.ValidAddress) == 0 {
		return nil, errors.New("未查询到有效收货地址，请前往 app 添加或检查填写的 cookie 是否正确！")
	}

	result := make(map[string]AddressItem)
	for _, v := range addressResult.Data.ValidAddress {
		result[v.Location.Address] = v
	}
	return result, nil
}
