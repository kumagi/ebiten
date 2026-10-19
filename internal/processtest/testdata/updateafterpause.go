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

//go:build ignore

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	phase   int
	draws   int
	updates int
	width   int
	height  int
	err     error
}

func (g *Game) Update() error {
	switch g.phase {
	case 0:
		g.updates++
		g.phase = 1
		ebiten.SetTPS(0)
	case 2:
		if g.width != 384 || g.height != 288 {
			return fmt.Errorf("resize Update at %dx%d, want 384x288", g.width, g.height)
		}
		g.updates++
		g.phase = 3
		g.draws = 0
	case 4:
		if g.err != nil {
			return g.err
		}
		return ebiten.Termination
	default:
		return fmt.Errorf("unexpected Update while TPS is zero: phase %d", g.phase)
	}
	return nil
}

func (g *Game) Draw(*ebiten.Image) {
	g.draws++
	switch g.phase {
	case 0:
		g.err = fmt.Errorf("Draw before the first Update")
	case 1:
		if g.draws < 2 {
			return
		}
		g.phase = 2
		g.draws = 0
		ebiten.SetWindowSize(384, 288)
		return
	case 2:
		if g.draws < 120 {
			return
		}
		g.err = fmt.Errorf("resize was not preceded by an Update")
	case 3:
		if g.draws < 2 {
			return
		}
		if g.updates != 2 {
			g.err = fmt.Errorf("Updates: got %d, want 2", g.updates)
		}
	default:
		return
	}
	g.phase = 4
	ebiten.SetTPS(ebiten.SyncWithFPS)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.width = width
	g.height = height
	return 64, 64
}

func main() {
	ebiten.SetWindowVisible(false)
	ebiten.SetWindowSize(320, 240)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
