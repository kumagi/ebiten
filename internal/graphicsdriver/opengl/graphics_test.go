// Copyright 2018 The Ebiten Authors
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

//go:build !playstation5

package opengl_test

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/internal/color"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl/gl"
	"testing"
)

type bufferContext struct {
	gl.Context
	created        int
	failAt         int
	bound          map[uint32]uint32
	live           map[uint32]bool
	invalidBinding bool
	uploads        int
}

func (c *bufferContext) CreateVertexArray() uint32                                   { return 1 }
func (c *bufferContext) BindVertexArray(uint32)                                      {}
func (c *bufferContext) EnableVertexAttribArray(uint32)                              {}
func (c *bufferContext) VertexAttribPointer(uint32, int32, uint32, bool, int32, int) {}
func (c *bufferContext) CreateBuffer() uint32 {
	c.created++
	if c.created == c.failAt {
		return 0
	}
	id := uint32(c.created)
	c.live[id] = true
	return id
}
func (c *bufferContext) BindBuffer(target, id uint32) {
	if id == 0 {
		c.invalidBinding = true
	}
	c.bound[target] = id
}
func (c *bufferContext) BufferInit(uint32, int, uint32) {}
func (c *bufferContext) DeleteBuffer(id uint32)         { delete(c.live, id) }
func (c *bufferContext) BufferSubData(target uint32, _ int, _ []byte) {
	c.uploads++
	if !c.live[c.bound[target]] {
		c.invalidBinding = true
	}
}

func TestSetVerticesBufferCreationFailure(t *testing.T) {
	for _, growing := range []bool{false, true} {
		for _, target := range []uint32{gl.ARRAY_BUFFER, gl.ELEMENT_ARRAY_BUFFER} {
			t.Run(fmt.Sprintf("growing=%t/target=%d", growing, target), func(t *testing.T) {
				c := &bufferContext{bound: map[uint32]uint32{}, live: map[uint32]bool{}}
				g := opengl.NewGraphicsForTesting(c, color.ColorSpaceSRGB)
				if growing {
					if err := g.SetVertices(make([]float32, 8), []uint32{0, 1, 2}); err != nil {
						t.Fatal(err)
					}
				}
				old := c.bound[target]
				c.failAt = c.created + 1
				if target == gl.ELEMENT_ARRAY_BUFFER {
					c.failAt++
				}
				c.uploads = 0
				if err := g.SetVertices(make([]float32, 16), []uint32{0, 1, 2, 0, 2, 3}); err == nil {
					t.Error("missing buffer creation error")
				}
				if c.invalidBinding || c.uploads != 0 {
					t.Errorf("failure: invalid binding=%v uploads=%d", c.invalidBinding, c.uploads)
				}
				if growing && (!c.live[old] || c.bound[target] != old) {
					t.Error("failed replacement discarded the old buffer")
				}
				c.failAt = 0
				c.uploads = 0
				if err := g.SetVertices(make([]float32, 16), []uint32{0, 1, 2, 0, 2, 3}); err != nil {
					t.Error(err)
				}
				if c.invalidBinding || c.uploads != 2 {
					t.Errorf("retry: invalid binding=%v uploads=%d", c.invalidBinding, c.uploads)
				}
			})
		}
	}
}
