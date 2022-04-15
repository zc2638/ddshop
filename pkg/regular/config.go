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

type Config struct {
	SuccessInterval int      `json:"success_interval"` // 执行成功 再次执行的间隔时间(ms), -1为停止继续执行
	FailInterval    int      `json:"fail_interval"`    // 执行失败 再次执行的间隔时间(ms), -1为停止继续执行
	Periods         []Period `json:"periods"`          // 执行周期
}

type Period struct {
	Start string `json:"start"` // 开始时间, 00:00
	End   string `json:"end"`   // 结束时间, 23:59

	startHour   int
	startMinute int
	endHour     int
	endMinute   int
}
