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

package app

import (
	"context"
	"errors"
	"os"

	"github.com/pkgms/go/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/zc2638/ddshop/asserts"
	"github.com/zc2638/ddshop/core/ddmc"
	"github.com/zc2638/ddshop/pkg/notice"
	"github.com/zc2638/ddshop/pkg/regular"
)

type Config struct {
	Bark     notice.BarkConfig     `json:"bark"`
	PushPlus notice.PushPlusConfig `json:"push_plus"`
	Regular  regular.Config        `json:"regular"`
	DDMC     ddmc.Config           `json:"ddmc"`
}

type Option struct {
	ConfigPath string
	Cookie     string
	BarkKey    string
	PayType    string
	Interval   int64
}

func (o *Option) Config() *Config {
	return &Config{
		Regular: regular.Config{
			SuccessInterval: 100,
			FailInterval:    100,
		},
		DDMC: ddmc.Config{
			Cookie:     o.Cookie,
			PayType:    o.PayType,
			Interval:   o.Interval,
			RetryCount: 100,
		},
		Bark: notice.BarkConfig{
			Key: o.BarkKey,
		},
	}
}

func NewRootCommand() *cobra.Command {
	opt := &Option{}
	cmd := &cobra.Command{
		Use:          "ddshop",
		Short:        "叮咚买菜自动抢购下单程序",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := opt.Config()
			if opt.ConfigPath != "" {
				viper.SetConfigType("yaml")
				if err := server.ParseConfigWithEnv(opt.ConfigPath, cfg, "DDSHOP"); err != nil {
					return err
				}
			} else {
				logrus.Warning("未设置配置文件，使用参数解析")
			}
			if cfg.DDMC.Cookie == "" {
				return errors.New("请输入用户Cookie.\n你可以执行此命令 `ddshop --cookie xxx` 或者 `DDSHOP_COOKIE=xxx ddshop`")
			}

			bark := notice.NewBark(&cfg.Bark)
			pushPlus := notice.NewPushPlus(&cfg.PushPlus)
			music := notice.NewMusic(asserts.NoticeMP3, 180)
			noticeIns := notice.New(notice.NewLog(), bark, pushPlus, music)

			session, err := ddmc.NewSession(&cfg.DDMC, noticeIns)
			if err != nil {
				return err
			}

			engine, err := regular.New(&cfg.Regular)
			if err != nil {
				return err
			}
			return engine.Start(context.Background(), session)
		},
	}

	configEnv := os.Getenv("DDSHOP_CONFIG")
	cookieEnv := os.Getenv("DDSHOP_COOKIE")
	barkKeyEnv := os.Getenv("DDSHOP_BARKKEY")
	payTypeEnv := os.Getenv("DDSHOP_PAYTYPE")
	cmd.Flags().StringVarP(&opt.ConfigPath, "config", "c", configEnv, "设置配置文件路径")
	cmd.Flags().StringVar(&opt.Cookie, "cookie", cookieEnv, "设置用户个人cookie")
	cmd.Flags().StringVar(&opt.BarkKey, "bark-key", barkKeyEnv, "设置bark的通知key")
	cmd.Flags().StringVar(&opt.PayType, "pay-type", payTypeEnv, "设置支付方式，支付宝、微信、alipay、wechat")
	cmd.Flags().Int64Var(&opt.Interval, "interval", 100, "设置请求间隔时间(ms)，默认为100")
	return cmd
}
