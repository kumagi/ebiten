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

package opengl_test

import (
	"syscall/js"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/color"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
)

func TestNewGraphicsInvalidColorSpace(t *testing.T) {
	calls := 0
	getContext := js.FuncOf(func(_ js.Value, _ []js.Value) any {
		calls++
		return js.Global().Get("Object").New()
	})
	defer getContext.Release()
	canvas := js.Global().Get("Object").New()
	canvas.Set("getContext", getContext)
	if got, err := opengl.NewGraphics(canvas, color.ColorSpace(99)); got != nil || err == nil {
		t.Errorf("NewGraphics = %v, %v; want nil and an error", got, err)
	}
	if calls != 0 {
		t.Errorf("getContext called %d times for an invalid color space", calls)
	}
}
