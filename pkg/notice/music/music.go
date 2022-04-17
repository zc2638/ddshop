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

package music

import (
	"bytes"
	"runtime"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

func NewMP3(b []byte, sec int) (*MP3, error) {
	decoder, err := mp3.NewDecoder(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return &MP3{
		decoder: decoder,
		sec:     sec,
	}, nil
}

type MP3 struct {
	decoder *mp3.Decoder
	sec     int
}

func (m *MP3) Play() error {
	c, _, err := oto.NewContext(m.decoder.SampleRate(), 2, 2)
	if err != nil {
		return err
	}
	player := c.NewPlayer(m.decoder)
	player.Play()

	// 异步放歌，需要等待
	time.Sleep(time.Duration(m.sec) * time.Second)
	runtime.KeepAlive(player)
	return nil
}

func (m *MP3) Name() string {
	return "Music"
}

func (m *MP3) Send(_, _ string) error {
	return m.Play()
}
