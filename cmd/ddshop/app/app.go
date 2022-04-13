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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/zc2638/ddshop/asserts"
	"github.com/zc2638/ddshop/core"
	"github.com/zc2638/ddshop/pkg/notice"
)

type Option struct {
	Cookie   string
	BarkKey  string
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
					logrus.Info("抢菜成功，请尽快支付!")
				})

				go func() {
					if opt.BarkKey == "" {
						return
					}
					ins := notice.NewBark(opt.BarkKey)
					if err := ins.Send("抢菜成功", "叮咚买菜 抢菜成功，请尽快支付！"); err != nil {
						logrus.Warningf("Bark消息通知失败: %v", err)
					}
				}()

				if err := asserts.Play(); err != nil {
					logrus.Warningf("播放成功提示音乐失败: %v", err)
				}
				// 异步放歌，歌曲有3分钟
				time.Sleep(3 * time.Minute)
				return nil
			}
		},
	}

	cookieEnv := os.Getenv("DDSHOP_COOKIE")
	barkKeyEnv := os.Getenv("DDSHOP_BARKKEY")
	cmd.Flags().StringVar(&opt.Cookie, "cookie", cookieEnv, "设置用户个人cookie")
	cmd.Flags().StringVar(&opt.BarkKey, "bark-key", barkKeyEnv, "设置bark的通知key")
	cmd.Flags().Int64Var(&opt.Interval, "interval", 500, "设置请求间隔时间(ms)，默认为100")
	return cmd
}

func Start(session *core.Session) error {
	logrus.Info("=====> 获取购物车中有效商品")

	cartData, err := session.GetCart()
	if err != nil {
		return err
	}
	if err := session.CartAllCheck(); err != nil {
		return fmt.Errorf("全选购物车商品失败: %v", err)
	}

	products := cartData["products"].([]map[string]interface{})
	for k, v := range products {
		logrus.Infof("[%v] %s 数量：%v 总价：%s", k, v["product_name"], v["count"], v["total_price"])
	}

	for {
		logrus.Info("=====> 获取可预约时间")
		multiReserveTime, err := session.GetMultiReserveTime(products)
		if err != nil {
			return fmt.Errorf("获取可预约时间失败: %v", err)
		}
		if len(multiReserveTime) == 0 {
			return core.ErrorNoReserveTime
		}
		logrus.Infof("发现可用的配送时段!")

		logrus.Info("=====> 生成订单信息")
		checkOrderData, err := session.CheckOrder(cartData, multiReserveTime)
		if err != nil {
			return fmt.Errorf("检查订单失败: %v", err)
		}
		logrus.Infof("订单总金额：%v\n", checkOrderData["price"])

		var wg errgroup.Group
		for _, reserveTime := range multiReserveTime {
			sess := session.Clone()
			sess.SetReserve(reserveTime)
			wg.Go(func() error {
				startTime := time.Unix(int64(sess.Reserve.StartTimestamp), 0).Format("2006/01/02 15:04:05")
				endTime := time.Unix(int64(sess.Reserve.EndTimestamp), 0).Format("2006/01/02 15:04:05")
				timeRange := startTime + "——" + endTime
				logrus.Infof("=====> 提交订单中, 预约时间段(%s)", timeRange)
				if err := sess.CreateOrder(context.Background(), cartData, checkOrderData); err != nil {
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
