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

package remote

import (
	"fmt"
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/color"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	"github.com/hajimehoshi/ebiten/v2/internal/vmprotocol"
)

type pixelTestHost struct{ sends, reads int }

func (h *pixelTestHost) SendGraphicsCommands([]vmprotocol.GraphicsCommand) error {
	h.sends++
	return nil
}
func (h *pixelTestHost) ReadPixelsAsync(graphicsdriver.ImageID, []graphicsdriver.PixelsArgs) (graphicsdriver.PixelReadback, error) {
	h.reads++
	return completedReadback{}, nil
}
func (*pixelTestHost) MaxImageSize() (int, error)            { return 32, nil }
func (*pixelTestHost) ColorSpace() (color.ColorSpace, error) { return color.ColorSpaceSRGB, nil }

func TestReadPixelsBufferLength(t *testing.T) {
	for _, hosted := range []bool{false, true} {
		for _, size := range []int{0, 3, 4, 5} {
			t.Run(fmt.Sprintf("host=%v/bytes=%d", hosted, size), func(t *testing.T) {
				g := NewGraphics()
				h := &pixelTestHost{}
				if hosted {
					g.SetHost(h)
				}
				g.commands = append(g.commands, vmprotocol.GraphicsCommand{Kind: vmprotocol.GraphicsCommandKindInitialize})
				img := &Image{id: 1, graphics: g}
				r, err := img.ReadPixelsAsync([]graphicsdriver.PixelsArgs{{Pixels: make([]byte, size), Region: image.Rect(0, 0, 1, 1)}})
				if size != 4 {
					if err == nil || r != nil {
						t.Errorf("invalid length accepted: %v,%v", r, err)
					}
					if h.sends != 0 || h.reads != 0 || len(g.commands) != 1 {
						t.Errorf("invalid request changed pending work: sends=%d,reads=%d,commands=%d", h.sends, h.reads, len(g.commands))
					}
				} else {
					if err != nil || r == nil {
						t.Fatalf("valid read: %v,%v", r, err)
					}
					r.Dispose()
					if hosted && (h.sends != 1 || h.reads != 1) {
						t.Errorf("valid request not forwarded: %v", h)
					}
				}
			})
		}
	}
}

func TestWritePixelsBufferLength(t *testing.T) {
	for _, size := range []int{0, 3, 4, 5} {
		t.Run(fmt.Sprintf("bytes=%d", size), func(t *testing.T) {
			g := NewGraphics()
			img := &Image{id: 1, graphics: g}
			valid := graphicsdriver.PixelsArgs{Pixels: []byte{1, 2, 3, 4}, Region: image.Rect(0, 0, 1, 1)}
			if err := img.WritePixels([]graphicsdriver.PixelsArgs{valid}); err != nil {
				t.Fatal(err)
			}
			err := img.WritePixels([]graphicsdriver.PixelsArgs{valid, {Pixels: make([]byte, size), Region: image.Rect(1, 0, 2, 1)}})
			if size != 4 {
				if err == nil {
					t.Error("invalid length accepted")
				}
				if len(g.commands) != 1 || len(g.regionsBuf) != 1 || len(g.pixelsBuf) != 4 || len(g.pixelsListBuf) != 1 {
					t.Errorf("rejected batch changed buffered state: commands=%d,regions=%d,pixels=%d,lists=%d", len(g.commands), len(g.regionsBuf), len(g.pixelsBuf), len(g.pixelsListBuf))
				}
			} else if err != nil || len(g.commands) != 2 || len(g.regionsBuf) != 3 || len(g.pixelsBuf) != 12 {
				t.Errorf("valid batch: %v,commands=%d,regions=%d,pixels=%d", err, len(g.commands), len(g.regionsBuf), len(g.pixelsBuf))
			}
		})
	}
}
