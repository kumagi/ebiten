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

package ebiten_test

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestTransformationStringValues(t *testing.T) {
	for _, transformed := range []bool{false, true} {
		var g ebiten.GeoM
		var s ebiten.ColorScale
		var m ebiten.ColorM
		if transformed {
			g.Scale(1.5, 0.5)
			g.Translate(2, 3)
			s.Scale(1.5, 0.5, 0.25, 0.75)
			m.Scale(1.5, 0.5, 0.25, 0.75)
		}
		for _, tc := range []struct {
			name    string
			value   any
			pointer fmt.Stringer
		}{
			{"GeoM", g, &g}, {"ColorScale", s, &s}, {"ColorM", m, &m},
		} {
			t.Run(fmt.Sprintf("%s/transformed=%v", tc.name, transformed), func(t *testing.T) {
				if _, ok := tc.value.(fmt.Stringer); !ok {
					t.Error("value does not implement fmt.Stringer")
				}
				if got, want := fmt.Sprint(tc.value), tc.pointer.String(); got != want {
					t.Errorf("value formatting = %q, want %q", got, want)
				}
				if got, want := fmt.Sprint(tc.pointer), tc.pointer.String(); got != want {
					t.Errorf("pointer formatting = %q, want %q", got, want)
				}
			})
		}
	}
	for _, pointer := range []any{(*ebiten.GeoM)(nil), (*ebiten.ColorScale)(nil), (*ebiten.ColorM)(nil)} {
		if got := fmt.Sprint(pointer); got != "<nil>" {
			t.Errorf("nil pointer formatting = %q, want <nil>", got)
		}
	}
}
