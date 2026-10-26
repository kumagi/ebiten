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

package ui_test

import (
	"github.com/hajimehoshi/ebiten/v2/internal/ui"
	"syscall/js"
	"testing"
)

func TestCanvasResizeKeepsUnchangedBuffer(t *testing.T) {
	document := js.Global().Get("document")
	if !document.Truthy() {
		t.Skip("requires a browser canvas")
	}
	body := document.Get("body")
	style := body.Get("style")
	oldWidth, oldHeight := style.Get("width").String(), style.Get("height").String()
	defer func() {
		style.Set("width", oldWidth)
		style.Set("height", oldHeight)
	}()
	style.Set("width", "160px")
	style.Set("height", "120px")
	c := document.Call("createElement", "canvas")
	ui.UpdateCanvasSizeForTest(c)
	ctx := c.Call("getContext", "2d")
	ctx.Set("fillStyle", "red")
	ctx.Call("fillRect", 0, 0, 1, 1)
	ui.UpdateCanvasSizeForTest(c)
	if alpha := ctx.Call("getImageData", 0, 0, 1, 1).Get("data").Index(3).Int(); alpha != 255 {
		t.Errorf("unchanged size cleared the drawing buffer: alpha=%d", alpha)
	}
	w := c.Get("width").Int()
	style.Set("width", "200px")
	ui.UpdateCanvasSizeForTest(c)
	if c.Get("width").Int() <= w {
		t.Error("canvas width did not follow the body resize")
	}
	if alpha := ctx.Call("getImageData", 0, 0, 1, 1).Get("data").Index(3).Int(); alpha != 0 {
		t.Errorf("changed size did not reset the drawing buffer: alpha=%d", alpha)
	}
}
