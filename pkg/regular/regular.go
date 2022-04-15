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

package regular

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type TaskInterface interface {
	Run(ctx context.Context) error
}

func New(cfg *Config) (*Engine, error) {
	for k, v := range cfg.Periods {
		start, err := time.Parse("15:04", v.Start)
		if err != nil {
			return nil, fmt.Errorf("解析时间段 %d 开始时间(%s)失败: %v", k, v.Start, err)
		}
		end, err := time.Parse("15:04", v.End)
		if err != nil {
			return nil, fmt.Errorf("解析时间段 %d 结束时间(%s)失败: %v", k, v.Start, err)
		}
		cfg.Periods[k].startHour = start.Hour()
		cfg.Periods[k].startMinute = start.Minute()
		cfg.Periods[k].endHour = end.Hour()
		cfg.Periods[k].endMinute = end.Minute()
	}
	return &Engine{cfg: cfg}, nil
}

type Engine struct {
	cfg    *Config
	cancel context.CancelFunc
}

func (e *Engine) Start(ctx context.Context, task TaskInterface) error {
	for {
		second := time.Now().Second()
		if second == 0 {
			break
		}
		sleepInterval := 60 - second
		logrus.Warningf("当前秒数不为 0，需等待 %ds 后开启自动助手", sleepInterval)
		time.Sleep(time.Duration(sleepInterval) * time.Second)
	}

	if len(e.cfg.Periods) == 0 {
		return e.run(ctx, task)
	}

	currentStartHour, currentStartMinute := -1, -1
	ticker := time.NewTicker(time.Minute)
	for {
		if e.cancel == nil {
			logrus.Infof("开始任务侦查")
		}
		now := time.Now()
		hour := now.Hour()
		minute := now.Minute()

		for _, v := range e.cfg.Periods {
			if currentStartHour > -1 && (currentStartHour != v.startHour || currentStartMinute != v.startMinute) {
				continue
			}

			start, end := false, false
			if v.startHour < hour {
				start = true
			}
			if v.startHour == hour && v.startMinute <= minute {
				start = true
			}
			if v.endHour < hour {
				end = true
			}
			if v.endHour == hour && v.endMinute <= minute {
				end = true
			}

			if start && !end && currentStartHour != v.startHour {
				currentStartHour = v.startHour
				currentStartMinute = v.startMinute

				ctx, e.cancel = context.WithCancel(ctx)
				go func() {
					if err := e.run(ctx, task); err != nil {
						logrus.Errorf("执行结束: %v", err)
					}
					logrus.Infof("当前时间段执行结束，请等待下个时间段")
				}()
				break
			}
			if start && end && e.cancel != nil {
				e.cancel()
				e.cancel = nil
			}
		}
		<-ticker.C
	}
}

func (e *Engine) run(ctx context.Context, task TaskInterface) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := task.Run(ctx); err != nil {
			if e.cfg.FailInterval < 0 {
				return err
			}
			logrus.Errorf("将在 %dms 后继续执行，执行出错: %v", e.cfg.FailInterval, err)
			fmt.Println()
			time.Sleep(time.Duration(e.cfg.FailInterval) * time.Millisecond)
			continue
		}

		if e.cfg.SuccessInterval < 0 {
			return nil
		}
		logrus.Infof("将在 %dms 后继续执行，执行成功", e.cfg.SuccessInterval)
		time.Sleep(time.Duration(e.cfg.SuccessInterval) * time.Millisecond)
	}
}
