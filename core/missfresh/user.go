package missfresh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type AddressDefault struct {
	Data struct {
		ReceiveAddressInfo struct {
			Id            int    `json:"id"`
			BusinessType  int    `json:"businessType"`
			LatLng        string `json:"latLng"`
			AreaCode      string `json:"areaCode"`
			City          string `json:"city"`
			Province      string `json:"province"`
			Address1      string `json:"address1"`
			Address2      string `json:"address2"`
			FullAddress   string `json:"fullAddress"`
			AddressDetail string `json:"addressDetail"`
			Name          string `json:"name"`
			PhoneNumber   string `json:"phoneNumber"`
			Gender        int    `json:"gender"`
			Tag           string `json:"tag"`
			PoiId         string `json:"poiId"`
			IsOpen        int    `json:"isOpen"`
			Optimal       bool   `json:"optimal"`
			Latlng        string `json:"lat_lng"`
		} `json:"receiveAddressInfo"`
		SwitchCity bool   `json:"switchCity"`
		CityName   string `json:"cityName"`
		Tips       string `json:"tips"`
	} `json:"data"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}

type DefaultAddressInfo struct {
	AddressID     int    `json:"address_id"`
	Username      string `json:"username"`
	AreaCode      string `json:"area_code"`
	FullAddress   string `json:"full_address"`
	AddressDetail string `json:"address_detail"`
	PoiId         string `json:"poi_id"`
	Lat           string `json:"lat"`
	Lng           string `json:"lng"`
}

func (s *Session) GetDefaultAddressView(ctx context.Context) error {
	u, err := url.Parse(HOST + User_API)
	if err != nil {
		return fmt.Errorf("userview url parse failed: %v", err)
	}

	params := s.buildURLParams(false)
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	marshal, err := json.Marshal(s.CreateBodyParam())
	if err != nil {
		return err
	}
	req.SetHeader("x-region", `{"address_code":"","station_code":"","delivery_type":"","bigWarehouse":"","type":null}`)
	req.SetBody(marshal)
	resp, err := s.execute(ctx, req, http.MethodPost, urlPath, maxRetryCount)
	if err != nil {
		return err
	}

	var userResult AddressDefault
	if err := json.Unmarshal(resp.Body(), &userResult); err != nil {
		return fmt.Errorf("parse response failed: %v, body: %v", err, resp.String())
	}

	arrs := strings.Split(userResult.Data.ReceiveAddressInfo.Latlng, ",")
	if len(arrs) != 2 {
		return fmt.Errorf("parse user view lat and lng failed")
	}

	res := &DefaultAddressInfo{
		Username:      userResult.Data.ReceiveAddressInfo.Name,
		AddressID:     userResult.Data.ReceiveAddressInfo.Id,
		AreaCode:      userResult.Data.ReceiveAddressInfo.AreaCode,
		FullAddress:   userResult.Data.ReceiveAddressInfo.FullAddress,
		AddressDetail: userResult.Data.ReceiveAddressInfo.AddressDetail,
		PoiId:         userResult.Data.ReceiveAddressInfo.PoiId,
		Lat:           arrs[0],
		Lng:           arrs[1],
	}
	s.DefaultAddressInfo = res
	logrus.Infof("获取默认地址信息成功, id: %d, name: %s, Address: %s", s.DefaultAddressInfo.AddressID, s.DefaultAddressInfo.Username, s.DefaultAddressInfo.FullAddress)
	return nil
}
