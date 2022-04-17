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

package notice

import (
	"github.com/sirupsen/logrus"

	"github.com/zc2638/ddshop/pkg/util"
)

func NewLog() Engine {
	return &log{}
}

type log struct{}

func (l *log) Name() string {
	return "Log"
}

func (l *log) Send(title, body string) error {
	util.LoopRun(10, func() {
		logrus.Infof("%s: %s", title, body)
	})
	return nil
}
