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
	"errors"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl/gl"
	"testing"
)

type initContext struct {
	gl.Context
	err   error
	loads int
	calls int
}

func (c *initContext) LoadFunctions() error                             { c.loads++; return c.err }
func (c *initContext) Enable(uint32)                                    { c.calls++ }
func (c *initContext) BlendEquationSeparate(uint32, uint32)             { c.calls++ }
func (c *initContext) BlendFuncSeparate(uint32, uint32, uint32, uint32) { c.calls++ }
func (c *initContext) GetInteger(uint32) int                            { c.calls++; return 0 }

func TestResetPreservesLoadError(t *testing.T) {
	loadError := errors.New("loading functions failed")
	for _, err := range []error{nil, loadError} {
		ctx := &initContext{err: err}
		reset, _ := opengl.NewContextForTesting(ctx)
		for i := range 3 {
			if got := reset(); got != err {
				t.Errorf("reset %d: error = %v, want %v", i, got, err)
			}
			if ctx.loads != 1 {
				t.Errorf("reset %d: LoadFunctions called %d times", i, ctx.loads)
			}
			if err != nil && ctx.calls != 0 {
				t.Errorf("reset %d: called GL methods %d times after load error", i, ctx.calls)
			}
		}
		if err == nil && ctx.calls == 0 {
			t.Error("successful reset did not call GL methods")
		}
	}
}
