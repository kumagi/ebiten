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

package main

import (
	"reflect"
	"testing"
)

func TestMoveCollision(t *testing.T) {
	for _, tc := range []struct {
		name      string
		body      []Position
		direction int
	}{
		{"left", []Position{{0, 5}}, dirLeft},
		{"right", []Position{{xGridCountInScreen - 1, 5}}, dirRight},
		{"top", []Position{{5, 0}}, dirUp},
		{"bottom", []Position{{5, yGridCountInScreen - 1}}, dirDown},
		{"self", []Position{{2, 2}, {2, 3}, {3, 3}, {3, 2}, {3, 1}}, dirRight},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := newGame()
			g.snakeBody, g.moveDirection, g.apple = tc.body, tc.direction, Position{0, 0}
			g.score, g.bestScore = 2, 9
			if err := g.Update(); err != nil {
				t.Fatal(err)
			}
			if g.moveDirection != dirNone || len(g.snakeBody) != 1 || g.score != 0 || g.bestScore != 9 || g.collidesWithWall() {
				t.Errorf("collision did not reset immediately: %+v", g)
			}
		})
	}
}

func TestMoveIntoVacatedTail(t *testing.T) {
	g := newGame()
	g.snakeBody = []Position{{2, 2}, {2, 3}, {3, 3}, {3, 2}}
	g.moveDirection, g.apple = dirRight, Position{0, 0}
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	if want := ([]Position{{3, 2}, {2, 2}, {2, 3}, {3, 3}}); !reflect.DeepEqual(g.snakeBody, want) {
		t.Errorf("body = %v, want %v", g.snakeBody, want)
	}
}

func TestMoveEatsAppleAndKeepsTail(t *testing.T) {
	g := newGame()
	g.snakeBody = []Position{{2, 2}, {1, 2}}
	g.moveDirection, g.apple, g.bestScore = dirRight, Position{3, 2}, 9
	if err := g.Update(); err != nil {
		t.Fatal(err)
	}
	if want := ([]Position{{3, 2}, {2, 2}, {1, 2}}); !reflect.DeepEqual(g.snakeBody, want) || g.score != 1 || g.bestScore != 9 {
		t.Errorf("body = %v, score/best = %d/%d; want %v, 1/9", g.snakeBody, g.score, g.bestScore, want)
	}
}
