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

package core

type Error string

func (e Error) Error() string {
	return string(e)
}

const ErrorComplete = Error("完成")

const (
	ErrorNoValidProduct     = Error("无有效商品")
	ErrorNoStock            = Error("部分商品已缺货")
	ErrorProductChange      = Error("商品信息有变化")
	ErrorInvalidReserveTime = Error("送达时间已失效")
	ErrorNoReserveTime      = Error("无可预约时间段")
)
