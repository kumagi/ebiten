// Copyright 2017 The Ebiten Authors
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

package audio_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type pausedLoopReader struct {
	*bytes.Reader
	emptyReads     int
	calls          int
	maxRead        int
	pauseAfterData int
}

func (r *pausedLoopReader) Read(p []byte) (int, error) {
	r.calls++
	if r.emptyReads > 0 {
		r.emptyReads--
		return 0, nil
	}
	if r.maxRead > 0 {
		p = p[:min(len(p), r.maxRead)]
	}
	n, err := r.Reader.Read(p)
	if n > 0 {
		r.emptyReads = r.pauseAfterData
	}
	return n, err
}

func TestInfiniteLoopEmptySourceReads(t *testing.T) {
	for _, f32 := range []bool{false, true} {
		frame := 4
		newLoop := audio.NewInfiniteLoop
		if f32 {
			frame = 8
			newLoop = audio.NewInfiniteLoopF32
		}
		for _, empty := range []int{0, 3, 99, 100} {
			t.Run(fmt.Sprintf("F32=%v/empty=%d", f32, empty), func(t *testing.T) {
				pcm := []byte{1, 2, 3, 4, 5, 6, 7, 8}
				src := &pausedLoopReader{Reader: bytes.NewReader(pcm), emptyReads: empty}
				loop := newLoop(src, 8)
				got := make([]byte, frame)
				n, err := loop.Read(got)
				if empty == 100 {
					if n != 0 || !errors.Is(err, io.ErrNoProgress) || src.calls != 100 {
						t.Fatalf("Read=%d/%v source calls=%d", n, err, src.calls)
					}
					n, err = loop.Read(got)
				}
				if n != frame || err != nil || !bytes.Equal(got, pcm[:frame]) {
					t.Fatalf("Read=%d/%v PCM=%v want=%v", n, err, got, pcm[:frame])
				}
			})
		}
		t.Run(fmt.Sprintf("F32=%v/partial progress", f32), func(t *testing.T) {
			pcm := []byte{1, 2, 3, 4, 5, 6, 7, 8}
			src := &pausedLoopReader{Reader: bytes.NewReader(pcm), maxRead: 1, pauseAfterData: 99}
			loop := newLoop(src, 8)
			got := make([]byte, frame)
			if n, err := loop.Read(got); n != frame || err != nil || !bytes.Equal(got, pcm[:frame]) {
				t.Fatalf("Read=%d/%v PCM=%v", n, err, got)
			}
		})
	}
}

func TestInfiniteLoopPartialFrameAfterNoProgress(t *testing.T) {
	for _, f32 := range []bool{false, true} {
		frame := 4
		newLoop := audio.NewInfiniteLoop
		if f32 {
			frame = 8
			newLoop = audio.NewInfiniteLoopF32
		}
		pcm := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		src := &pausedLoopReader{Reader: bytes.NewReader(pcm), maxRead: 1, pauseAfterData: 100}
		loop := newLoop(src, 8)
		got := make([]byte, frame)
		if n, err := loop.Read(nil); n != 0 || err != nil || src.calls != 0 {
			t.Fatalf("empty Read=%d/%v source calls=%d", n, err, src.calls)
		}
		if n, err := loop.Read(got); n != 0 || !errors.Is(err, io.ErrNoProgress) {
			t.Fatalf("partial stalled Read=%d/%v", n, err)
		}
		src.emptyReads = 0
		src.pauseAfterData = 0
		if n, err := loop.Read(got); n != frame || err != nil || !bytes.Equal(got, pcm[:frame]) {
			t.Fatalf("recovered Read=%d/%v PCM=%v", n, err, got)
		}
	}
}
