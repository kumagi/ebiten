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

package textinput_test

import (
	"math"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
)

func TestCompositionSelectionRange(t *testing.T) {
	for _, tc := range []struct {
		name, text                     string
		start, end, wantStart, wantEnd int
	}{
		{"valid", "abc", 1, 2, 1, 2},
		{"negative", "abc", -1, -1, 0, 0},
		{"reversed", "abc", 2, 1, 2, 2},
		{"past end", "abc", 100, 200, 3, 3},
		{"whole text", "abc", -1, 100, 0, 3},
		{"integer bounds", "abc", math.MinInt, math.MaxInt, 0, 3},
		{"empty", "", math.MaxInt, math.MinInt, 0, 0},
		{"unicode", "あ🚀b", 3, 7, 3, 7},
		{"byte offsets", "あ🚀b", 4, 5, 4, 5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			driver := textinput.NewComposerDriver("", "")
			defer func() {
				driver.Composer.OnComposition = nil
				driver.Composer.Cancel()
			}()
			var callbacks int
			driver.Composer.OnComposition = func(c *textinput.Composition) {
				callbacks++
				start, end := c.SelectionRangeInBytes()
				if c.Text() != tc.text || start != tc.wantStart || end != tc.wantEnd {
					t.Errorf("composition=%q range=[%d,%d), want %q [%d,%d)", c.Text(), start, end, tc.text, tc.wantStart, tc.wantEnd)
				}
			}
			if !driver.Send(textinput.TextInputState{
				Text:                             tc.text,
				CompositionSelectionStartInBytes: tc.start,
				CompositionSelectionEndInBytes:   tc.end,
			}) {
				t.Fatal("session did not accept composition")
			}
			if _, err := driver.Composer.Update(); err != nil {
				t.Fatal(err)
			}
			if callbacks != 1 {
				t.Errorf("callbacks=%d, want 1", callbacks)
			}
		})
	}
}
