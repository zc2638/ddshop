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

	"github.com/sirupsen/logrus"
)

type UserResult struct {
	Success bool     `json:"success"`
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type UserData struct {
	DoingRefundNum      int         `json:"doing_refund_num"`
	NoCommentOrderPoint int         `json:"no_comment_order_point"`
	NameNotice          string      `json:"name_notice"`
	NoPayOrderNum       int         `json:"no_pay_order_num"`
	DoingOrderNum       int         `json:"doing_order_num"`
	UserVip             UserVIP     `json:"user_vip"`
	UserSign            UserSign    `json:"user_sign"`
	NotOnionTip         int         `json:"not_onion_tip"`
	NoDrawCouponMoney   string      `json:"no_draw_coupon_money"`
	PointNum            int         `json:"point_num"`
	Balance             UserBalance `json:"balance"`
	UserInfo            UserInfo    `json:"user_info"`
	CouponNum           int         `json:"coupon_num"`
	NoCommentOrderNum   int         `json:"no_comment_order_num"`
}

type UserVIP struct {
	IsRenew                  int    `json:"is_renew"`
	VipSaveMoneyDescription  string `json:"vip_save_money_description"`
	VipDescription           string `json:"vip_description"`
	VipStatus                int    `json:"vip_status"`
	VipNotice                string `json:"vip_notice"`
	VipExpireTimeDescription string `json:"vip_expire_time_description"`
	VipUrl                   string `json:"vip_url"`
}

type UserSign struct {
	IsTodaySign bool   `json:"is_today_sign"`
	SignSeries  int    `json:"sign_series"`
	SignText    string `json:"sign_text"`
}

type UserBalance struct {
	SetFingerPayPassword int    `json:"set_finger_pay_password"`
	Balance              string `json:"balance"`
	SetPayPassword       int    `json:"set_pay_password"`
}

type UserInfo struct {
	Birthday       string `json:"birthday"`
	ShowInviteCode bool   `json:"show_invite_code"`
	NameInCheck    string `json:"name_in_check"`
	InviteCodeUrl  string `json:"invite_code_url"`
	Sex            int    `json:"sex"`
	Mobile         string `json:"mobile"`
	Avatar         string `json:"avatar"`
	ImUid          int    `json:"im_uid"`
	BindStatus     int    `json:"bind_status"`
	NameStatus     int    `json:"name_status"`
	NewRegister    bool   `json:"new_register"`
	ImSecret       string `json:"im_secret"`
	Name           string `json:"name"`
	Id             string `json:"id"`
	Introduction   string `json:"introduction"`
}

func (s *Session) GetUser() error {
	u, err := url.Parse("https://sunquan.api.ddxq.mobi/api/v1/user/detail/")
	if err != nil {
		return fmt.Errorf("user url parse failed: %v", err)
	}

	params := s.buildURLParams(false)
	u.RawQuery = params.Encode()
	urlPath := u.String()

	req := s.client.R()
	req.SetHeader("Host", "sunquan.api.ddxq.mobi")
	resp, err := s.execute(context.Background(), req, http.MethodGet, urlPath, maxRetryCount)
	if err != nil {
		return err
	}

	var userResult UserResult
	if err := json.Unmarshal(resp.Body(), &userResult); err != nil {
		return fmt.Errorf("parse response failed: %v, body: %v", err, resp.String())
	}

	s.UserID = userResult.Data.UserInfo.Id
	logrus.Infof("获取用户信息成功, id: %s, name: %s", s.UserID, userResult.Data.UserInfo.Name)
	return nil
}
