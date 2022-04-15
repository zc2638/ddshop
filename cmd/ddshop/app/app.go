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
	"github.com/spf13/viper"

	"github.com/zc2638/ddshop/core/missfresh"

	"github.com/zc2638/ddshop/asserts"

	"github.com/zc2638/ddshop/pkg/notice"

	"github.com/zc2638/ddshop/pkg/regular"

	"github.com/zc2638/ddshop/core/ddmc"

	"github.com/spf13/cobra"
)

type Config struct {
	Bark      notice.BarkConfig `json:"bark"`
	Regular   regular.Config    `json:"regular"`
	DDMC      ddmc.Config       `json:"ddmc"`
	Missfresh missfresh.Config  `json:"missfresh"`
}

type Option struct {
	ConfigPath string
	Cookie     string
	Token      string
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
			Cookie:   o.Cookie,
			Interval: o.Interval,
			PayType:  o.PayType,
		},
		Bark: notice.BarkConfig{
			Key: o.BarkKey,
		},
		Missfresh: missfresh.Config{
			Token:    o.Token,
			Interval: o.Interval,
			PayType:  o.PayType,
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
			if opt.ConfigPath == "" {
				return errors.New("未设置配置文件")
			}
			viper.SetConfigType("yaml")
			err := server.ParseConfigWithEnv(opt.ConfigPath, cfg, "DDSHOP")
			if err != nil {
				return err
			}

			bark := notice.NewBark(&cfg.Bark)
			music := notice.NewMusic(asserts.NoticeMP3, 180)
			noticeIns := notice.New(notice.NewLog(), bark, music)

			var session regular.TaskInterface
			if cfg.DDMC.Cookie != "" {
				session, err = ddmc.NewSession(&cfg.DDMC, noticeIns)
				if err != nil {
					return err
				}
			} else if cfg.Missfresh.Token != "" {
				session, err = missfresh.NewSession(&cfg.Missfresh, noticeIns)
				if err != nil {
					return err
				}
			} else {
				return errors.New("cookie和token均未设置，不执行抢购程序")
			}
			engine, err := regular.New(&cfg.Regular)
			if err != nil {
				return err
			}
			return engine.Start(context.Background(), session)
		},
	}

	configEnv := os.Getenv("DDSHOP_CONFIG")
	cmd.Flags().StringVarP(&opt.ConfigPath, "config", "c", configEnv, "设置配置文件路径")
	return cmd
}
