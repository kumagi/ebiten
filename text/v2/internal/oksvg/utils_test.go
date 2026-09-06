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

package oksvg_test

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2/internal/oksvg"
)

func TestParseSVGColorOutOfRange(t *testing.T) {
	for _, test := range []struct {
		in   string
		want color.NRGBA
	}{
		{in: "rgb(-1,0,0)", want: color.NRGBA{R: 0, G: 0, B: 0, A: 0xff}},
		{in: "rgb(0,-1,0)", want: color.NRGBA{R: 0, G: 0, B: 0, A: 0xff}},
		{in: "rgb(256,0,0)", want: color.NRGBA{R: 0xff, G: 0, B: 0, A: 0xff}},
		{in: "rgb(-1%,0%,0%)", want: color.NRGBA{R: 0, G: 0, B: 0, A: 0xff}},
		{in: "rgb(1000%,0%,0%)", want: color.NRGBA{R: 0xff, G: 0, B: 0, A: 0xff}},
		{in: "rgb(100%,0%,0%)", want: color.NRGBA{R: 0xff, G: 0, B: 0, A: 0xff}},
	} {
		got, err := oksvg.ParseSVGColor(test.in)
		if err != nil {
			t.Errorf("ParseSVGColor(%q) returned an error: %v", test.in, err)
			continue
		}
		if got != test.want {
			t.Errorf("ParseSVGColor(%q) = %v, want %v", test.in, got, test.want)
		}
	}
}
