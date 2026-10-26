// Copyright 2019 The Ebiten Authors
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

// Package mobile provides functions for mobile platforms (Android and iOS).
//
// This package is used when you use `ebitenmobile bind`.
//
// For usage, see https://ebitengine.org/en/documents/mobile.html.
package mobile

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// SetGame sets a mobile game.
//
// SetGame can be called at any time. Until SetGame or SetGameWithOptions is called, the game does not start.
//
// SetGame and SetGameWithOptions must be called only once in total. A subsequent call panics.
// SetGame also panics on platforms other than Android and iOS.
func SetGame(game ebiten.Game) {
	SetGameWithOptions(game, nil)
}

// SetGameWithOptions sets a mobile game with the specified options.
//
// SetGameWithOptions can be called at any time. Until SetGameWithOptions or SetGame is called, the game does not start.
//
// SetGameWithOptions and SetGame must be called only once in total. A subsequent call panics.
// SetGameWithOptions also panics on platforms other than Android and iOS.
func SetGameWithOptions(game ebiten.Game, options *ebiten.RunGameOptions) {
	setGame(game, options)
}
