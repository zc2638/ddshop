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
	"errors"
	"fmt"
	"sync"
)

type Engine interface {
	Name() string
	Send(title, body string) error
}

type Interface interface {
	Notice(title, body string) error
}

func New(engines ...Engine) Interface {
	es := make([]Engine, 0, len(engines))
	for _, e := range engines {
		if e == nil {
			continue
		}
		es = append(es, e)
	}
	return &notice{engines: es}
}

type notice struct {
	engines []Engine
}

func (n *notice) Notice(title, body string) error {
	var mux sync.Mutex
	var errStr string

	var wg sync.WaitGroup
	for _, v := range n.engines {
		wg.Add(1)

		go func(engine Engine) {
			defer wg.Done()

			if err := engine.Send(title, body); err != nil {
				mux.Lock()
				defer mux.Unlock()
				errStr += fmt.Sprintf("%s: %v\n", engine.Name(), err)
			}
		}(v)
	}
	wg.Wait()

	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}
