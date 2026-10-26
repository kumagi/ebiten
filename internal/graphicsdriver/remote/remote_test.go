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

package remote_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/color"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/remote"
)

type queryHost struct {
	remote.Host
	err        error
	colorCalls int
	sizeCalls  int
}

func (h *queryHost) ColorSpace() (color.ColorSpace, error) {
	h.colorCalls++
	return color.ColorSpaceSRGB, h.err
}

func (h *queryHost) MaxImageSize() (int, error) {
	h.sizeCalls++
	return 4096, h.err
}

func TestHostQueryCache(t *testing.T) {
	tests := []struct {
		name  string
		query func(*remote.Graphics) int
		calls func(*queryHost) int
		want  int
	}{
		{
			name:  "ColorSpace",
			query: func(g *remote.Graphics) int { return int(g.ColorSpace()) },
			calls: func(h *queryHost) int { return h.colorCalls },
			want:  int(color.ColorSpaceSRGB),
		},
		{
			name:  "MaxImageSize",
			query: func(g *remote.Graphics) int { return g.MaxImageSize() },
			calls: func(h *queryHost) int { return h.sizeCalls },
			want:  4096,
		},
	}
	for _, test := range tests {
		for _, failed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%v", test.name, failed), func(t *testing.T) {
				h := &queryHost{}
				if failed {
					h.err = errors.New("query failed")
				}
				g := remote.NewGraphics()
				g.SetHost(h)
				for range 2 {
					func() {
						if failed {
							defer func() {
								if _, ok := recover().(string); !ok {
									t.Error("failed query did not panic with a string")
								}
							}()
						}
						if got := test.query(g); !failed && got != test.want {
							t.Errorf("got %d, want %d", got, test.want)
						}
					}()
					h.err = nil
				}
				if got := test.calls(h); got != 1 {
					t.Errorf("host queries = %d, want 1", got)
				}
			})
		}
	}
}
