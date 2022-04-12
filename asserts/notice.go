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

package asserts

import (
	"bytes"
	_ "embed"
	"io"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

// noticeMP3 defines the content of notice.mp3
//go:embed notice.mp3
var noticeMP3 []byte

func Play() error {
	rc := io.NopCloser(bytes.NewReader(noticeMP3))
	streamer, format, err := mp3.Decode(rc)
	if err != nil {
		return err
	}
	defer streamer.Close()

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		return err
	}
	speaker.Play(streamer)
	return nil
}
