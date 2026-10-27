// Copyright 2026 The Ebitengine Authors
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

package graphicscommand

import (
	"errors"
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
)

type abortedReadbackImage struct {
	graphicsdriver.Image
	reads int
}

func (i *abortedReadbackImage) ReadPixelsAsync([]graphicsdriver.PixelsArgs) (graphicsdriver.PixelReadback, error) {
	i.reads++
	return nil, errors.New("unexpected read-back")
}

func TestAbortedReadPixelsCommandDoesNotExecute(t *testing.T) {
	for _, created := range []bool{false, true} {
		req := newReadPixelsRequest([]graphicsdriver.PixelsArgs{{Pixels: make([]byte, 4), Region: image.Rect(0, 0, 1, 1)}})
		img := &Image{}
		driverImage := &abortedReadbackImage{}
		if created {
			img.image = driverImage
		}
		c := &readPixelsAsyncCommand{img: img, req: req}
		q := &commandQueue{commands: []command{c}}
		want := errors.New("aborted")
		q.abortQueuedReadPixels(want)
		if err := c.Exec(q, nil, 0); err != nil {
			t.Fatal(err)
		}
		if done, err := req.result.Poll(); !done || err != want {
			t.Errorf("created=%v: result=(%v,%v)", created, done, err)
		}
		if driverImage.reads != 0 {
			t.Errorf("created=%v: aborted command started %d read-backs", created, driverImage.reads)
		}
	}
}
