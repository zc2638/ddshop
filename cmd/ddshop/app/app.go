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
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/zc2638/ddshop/asserts"

	"golang.org/x/sync/errgroup"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/zc2638/ddshop/core"
)

type Option struct {
	Cookie   string
	Interval int64
}

var (
	successCh = make(chan struct{}, 1)
	errCh     = make(chan error, 1)
)

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

			session := core.NewSession(opt.Cookie, opt.Interval)
			if err := session.GetUser(); err != nil {
				return fmt.Errorf("获取用户信息失败: %v", err)
			}
			if err := session.Choose(); err != nil {
				return err
			}

			go func() {
				for {
					if err := Start(session); err != nil {
						switch err {
						case core.ErrorNoValidProduct:
							logrus.Error("购物车中无有效商品，请先前往app添加或勾选！")
							errCh <- err
							return
						case core.ErrorNoReserveTime:
							sleepInterval := 3 + rand.Intn(6)
							logrus.Warningf("暂无可预约的时间，%d 秒后重试！", sleepInterval)
							time.Sleep(time.Duration(sleepInterval) * time.Second)
						default:
							logrus.Error(err)
						}
						fmt.Println()
					}
				}
			}()

			select {
			case err := <-errCh:
				return err
			case <-successCh:
				core.LoopRun(10, func() {
					logrus.Info("抢到菜了，请速去支付!")
				})

				if err := asserts.Play(); err != nil {
					logrus.Warningf("播放成功提示音乐失败: %v", err)
				}
				return nil
			}
		},
	}

	cookieEnv := os.Getenv("DDSHOP_COOKIE")
	cmd.Flags().StringVar(&opt.Cookie, "cookie", cookieEnv, "设置用户个人cookie")
	cmd.Flags().Int64Var(&opt.Interval, "interval", 100, "设置请求间隔时间(ms)，默认为100")
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

	if err := session.CartAllCheck(); err != nil {
		return fmt.Errorf("全选购车车商品失败: %v", err)
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

		var wg errgroup.Group
		for _, reserveTime := range multiReserveTime {
			sess := session.Clone()
			sess.UpdatePackageOrder(reserveTime)
			wg.Go(func() error {
				startTime := time.Unix(int64(sess.PackageOrder.PaymentOrder.ReservedTimeStart), 0).Format("2006/01/02 15:04:05")
				endTime := time.Unix(int64(sess.PackageOrder.PaymentOrder.ReservedTimeEnd), 0).Format("2006/01/02 15:04:05")
				timeRange := startTime + "——" + endTime
				logrus.Infof(">>> 提交订单中, 预约时间段(%s)", timeRange)
				if err := sess.CreateOrder(context.Background()); err != nil {
					logrus.Warningf("提交订单(%s)失败: %v", timeRange, err)
					return err
				}

				successCh <- struct{}{}
				return nil
			})
		}
		_ = wg.Wait()
		return nil
	}
}
