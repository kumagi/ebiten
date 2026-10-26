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
	"github.com/hajimehoshi/ebiten/v2/internal/color"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl/gl"
	"testing"
)

type vertexArrayContext struct {
	gl.Context
	fail           bool
	arrays         int
	invalidBinding bool
	uploads        int
}

func (c *vertexArrayContext) CreateVertexArray() uint32 {
	c.arrays++
	if c.fail {
		return 0
	}
	return 1
}
func (c *vertexArrayContext) BindVertexArray(id uint32) {
	if id == 0 {
		c.invalidBinding = true
	}
}
func (c *vertexArrayContext) CreateBuffer() uint32                                        { return 1 }
func (c *vertexArrayContext) BindBuffer(uint32, uint32)                                   {}
func (c *vertexArrayContext) BufferInit(uint32, int, uint32)                              {}
func (c *vertexArrayContext) BufferSubData(uint32, int, []byte)                           { c.uploads++ }
func (c *vertexArrayContext) EnableVertexAttribArray(uint32)                              {}
func (c *vertexArrayContext) VertexAttribPointer(uint32, int32, uint32, bool, int32, int) {}

func TestSetVerticesVertexArrayCreationFailure(t *testing.T) {
	c := &vertexArrayContext{fail: true}
	g := opengl.NewGraphicsForTesting(c, color.ColorSpaceSRGB)
	vertices := make([]float32, 8)
	indices := []uint32{0, 1, 2}
	if err := g.SetVertices(vertices, indices); err == nil {
		t.Error("missing vertex array creation error")
	}
	if c.invalidBinding || c.uploads != 0 {
		t.Errorf("failure: invalid vertex array binding=%v uploads=%d", c.invalidBinding, c.uploads)
	}
	c.fail = false
	if err := g.SetVertices(vertices, indices); err != nil {
		t.Error(err)
	}
	if c.invalidBinding || c.uploads != 2 {
		t.Errorf("retry: invalid vertex array binding=%v uploads=%d", c.invalidBinding, c.uploads)
	}
	if err := g.SetVertices(vertices, indices); err != nil {
		t.Error(err)
	}
	if c.arrays != 2 || c.uploads != 4 {
		t.Errorf("reuse: vertex array creations=%d uploads=%d", c.arrays, c.uploads)
	}
}
