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

package main

import "testing"

func TestScoreMatchesPipePass(t *testing.T) {
	// These are the collision dimensions used by hit, independent of sprite padding.
	const hitWidth, hitHeight = 30, 60
	w, h := gopherImage.Bounds().Dx(), gopherImage.Bounds().Dy()
	for n := 1; n <= 3; n++ {
		column := pipeStartOffsetX + n*pipeIntervalX
		boundary16 := (column*tileSize + pipeWidth - (w-hitWidth)/2) * 16
		for _, delta16 := range []int{-16, -1, 0, 1, 16} {
			g := &Game{mode: ModeGame, x16: boundary16 + delta16, y16: (5*tileSize - 1 - (h-hitHeight)/2) * 16, pipeTileYs: []int{5, 5, 5, 5}}
			want := n
			if delta16 < 0 {
				want--
			}
			if got := g.score(); got != want {
				t.Errorf("pipe=%d delta16=%d: score=%d want=%d", n, delta16, got, want)
			}
			if got := g.hit(); got != (delta16 < 0) {
				t.Errorf("pipe=%d delta16=%d: collision=%v", n, delta16, got)
			}
		}
	}
	for _, x16 := range []int{-500, -16, 0} {
		if got := (&Game{x16: x16}).score(); got != 0 {
			t.Errorf("initial x16=%d score=%d", x16, got)
		}
	}
}
