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
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/remote"
	"github.com/hajimehoshi/ebiten/v2/internal/vmprotocol"
)

type imageSizeHost struct {
	remote.Host
	commands []vmprotocol.GraphicsCommand
}

func (*imageSizeHost) MaxImageSize() (int, error) { return 32, nil }

func (h *imageSizeHost) SendGraphicsCommands(commands []vmprotocol.GraphicsCommand) error {
	h.commands = slices.Clone(commands)
	return nil
}

func TestImageSize(t *testing.T) {
	for _, tc := range []struct {
		name   string
		create func(*remote.Graphics, int, int) (graphicsdriver.Image, error)
		kind   vmprotocol.GraphicsCommandKind
	}{
		{
			name:   "image",
			create: (*remote.Graphics).NewImage,
			kind:   vmprotocol.GraphicsCommandKindNewImage,
		},
		{
			name:   "screen",
			create: (*remote.Graphics).NewScreenFramebufferImage,
			kind:   vmprotocol.GraphicsCommandKindNewScreenFramebufferImage,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, size := range [][2]int{{0, 1}, {1, 0}, {-1, 1}, {1, -1}, {33, 1}, {1, 33}} {
				host := &imageSizeHost{}
				g := remote.NewGraphics()
				g.SetHost(host)
				panicked := false
				func() {
					defer func() { panicked = recover() != nil }()
					if _, err := tc.create(g, size[0], size[1]); err != nil {
						t.Error(err)
					}
				}()
				if !panicked {
					t.Errorf("size %v must panic", size)
					continue
				}
				img, err := tc.create(g, 32, 1)
				if err != nil {
					t.Fatal(err)
				}
				if img.ID() != 1 {
					t.Errorf("image ID=%d, want 1 after rejected creation", img.ID())
				}
				if err := g.Flush(); err != nil {
					t.Fatal(err)
				}
				if len(host.commands) != 1 {
					t.Errorf("recorded commands=%d, want 1", len(host.commands))
					continue
				}
				got := host.commands[0]
				if got.Kind != tc.kind || got.Width != 32 || got.Height != 1 {
					t.Errorf("unexpected command: %+v", got)
				}
			}
		})
	}
}
