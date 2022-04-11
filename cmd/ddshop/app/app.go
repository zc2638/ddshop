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
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/zc2638/ddshop/core"
)

type Option struct {
	Cookie string
	Number int
}

func NewRootCommand() *cobra.Command {
	opt := &Option{}
	cmd := &cobra.Command{
		Use:          "ddshop",
		Short:        "Ding Dong grocery shopping automatic order program",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opt.Cookie == "" {
				return errors.New("请输入用户Cookie.\n你可以执行此命令 `ddshop --cookie xxx` 或者 `DDSHOP_COOKIE=xxx ddshop`")
			}

			session := core.NewSession(opt.Cookie, opt.Number)
			if err := session.GetUser(); err != nil {
				return fmt.Errorf("获取用户信息失败: %v", err)
			}
			if err := session.Choose(); err != nil {
				return err
			}

			for {
				if err := Start(session); err != nil {
					switch err {
					case core.ErrorNoValidProduct:
						logrus.Error("购物车中无有效商品，请先前往app添加或勾选！")
						return err
					case core.ErrorComplete:
						logrus.Info("抢购结束")
						return nil
					case core.ErrorNoReserveTime:
						sleepInterval := 3 + rand.Intn(6)
						logrus.Warningf("暂无可预约的时间，%d 秒后重试！", sleepInterval)
						time.Sleep(time.Duration(sleepInterval) * time.Second)
					default:
						logrus.Error(err)
					}
					println()
				}
			}
		},
	}

	cookieEnv := os.Getenv("DDSHOP_COOKIE")
	cmd.Flags().StringVar(&opt.Cookie, "cookie", cookieEnv, "设置用户个人cookie")
	cmd.Flags().IntVar(&opt.Number, "number", 1, "设置并行处理数量")
	return cmd
}

func Start(session *core.Session) error {
	logrus.Info(">>> 获取购物车中有效商品")
	if err := session.GetCart(); err != nil {
		return fmt.Errorf("检查购物车失败: %v", err)
	}
	if len(session.Cart.ProdList) == 0 {
		return core.ErrorNoValidProduct
	}

	for index, prod := range session.Cart.ProdList {
		logrus.Infof("[%v] %s 数量：%v 总价：%s", index, prod.ProductName, prod.Count, prod.TotalPrice)
	}
	session.Order.Products = session.Cart.ProdList

	for {
		logrus.Info(">>> 生成订单信息")
		if err := session.CheckOrder(); err != nil {
			return fmt.Errorf("检查订单失败: %v", err)
		}
		logrus.Infof("订单总金额：%v\n", session.Order.Price)

		session.GeneratePackageOrder()

		logrus.Info(">>> 获取可预约时间")
		multiReserveTime, err := session.GetMultiReserveTime()
		if err != nil {
			return fmt.Errorf("获取可预约时间失败: %v", err)
		}
		if len(multiReserveTime) == 0 {
			return core.ErrorNoReserveTime
		}
		logrus.Infof("发现可用的配送时段!")

		for _, reserveTime := range multiReserveTime {
			session.UpdatePackageOrder(reserveTime)

			for {
				logrus.Info(">>> 提交订单中")
				if err := session.CreateOrder(); err != nil {
					switch err {
					case core.ErrorInvalidReserveTime:
						// 选择的时间失效，使用下一个时间段
						logrus.Warningf("送达时间已失效, [%s]", reserveTime.SelectMsg)
						break
					case core.ErrorProductChange:
						return err
					default:
						logrus.Warningf("提交订单失败: %v", err)
						continue
					}
				}

				core.LoopRun(10, func() {
					logrus.Info("抢到菜了，请速去支付!")
				})
				return core.ErrorComplete
			}
		}
		return nil
	}
}
