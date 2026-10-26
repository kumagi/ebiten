// Copyright 2016 Hajime Hoshi
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
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl/gl"
	"testing"
)

type limitContext struct {
	gl.Context
	limit   int
	queries int
}

func (*limitContext) LoadFunctions() error                             { return nil }
func (*limitContext) Enable(uint32)                                    {}
func (*limitContext) BlendEquationSeparate(uint32, uint32)             {}
func (*limitContext) BlendFuncSeparate(uint32, uint32, uint32, uint32) {}
func (c *limitContext) GetInteger(name uint32) int {
	if name == gl.MAX_TEXTURE_SIZE {
		c.queries++
		return c.limit
	}
	if name == gl.FRAMEBUFFER_BINDING {
		return 0
	}
	panic("unexpected integer query")
}
func TestMaxTextureSizeAfterReset(t *testing.T) {
	ctx := &limitContext{}
	reset, maxTextureSize := opengl.NewContextForTesting(ctx)
	// Reset must invalidate the cache even when the replacement context has a smaller limit.
	for i, limit := range []int{16, 8, 32, 32} {
		ctx.limit = limit
		if err := reset(); err != nil {
			t.Fatal(err)
		}
		for range 2 {
			if got := maxTextureSize(); got != limit {
				t.Errorf("reset %d: limit = %d, want %d", i, got, limit)
			}
			if ctx.queries != i+1 {
				t.Errorf("reset %d: queries = %d, want %d", i, ctx.queries, i+1)
			}
		}
	}
}
