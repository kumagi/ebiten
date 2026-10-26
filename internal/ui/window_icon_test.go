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

//go:build !android && !ios && !js && !nintendosdk && !playstation5

package ui_test

import (
	"github.com/hajimehoshi/ebiten/v2/internal/ui"
	"image"
	"image/color"
	"testing"
)

type countingIcon struct {
	image.Image
	reads int
}

func (i *countingIcon) At(x, y int) color.Color {
	i.reads++
	return i.Image.At(x, y)
}

func TestPendingWindowIconConvertedOnce(t *testing.T) {
	img := &countingIcon{Image: image.NewRGBA(image.Rect(2, 3, 4, 6))}
	pending := ui.NewWindowIconImagesForTest([]image.Image{img})
	if img.reads != 0 {
		t.Fatal("icon pixels were read before game-loop processing")
	}
	for range 2 {
		converted := pending()
		if len(converted) != 1 || converted[0].Bounds() != img.Bounds() {
			t.Fatal("converted icon bounds changed")
		}
	}
	if got, want := img.reads, img.Bounds().Dx()*img.Bounds().Dy(); got != want {
		t.Errorf("pending icon pixels read %d times, want %d", got, want)
	}
}
