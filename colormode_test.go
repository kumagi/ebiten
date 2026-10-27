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

func TestInvalidPreferredColorMode(t *testing.T) {
	original := ebiten.PreferredColorMode()
	defer ebiten.SetPreferredColorMode(original)
	for _, invalid := range []ebiten.ColorMode{-1, 3, 99} {
		t.Run(fmt.Sprint(invalid), func(t *testing.T) {
			ebiten.SetPreferredColorMode(ebiten.ColorModeLight)
			func() {
				defer func() {
					if got, want := fmt.Sprint(recover()), fmt.Sprintf("ebiten: invalid color mode: %d", invalid); got != want {
						t.Errorf("panic = %q, want %q", got, want)
					}
				}()
				ebiten.SetPreferredColorMode(invalid)
			}()
			if got := ebiten.PreferredColorMode(); got != ebiten.ColorModeLight {
				t.Errorf("invalid mode changed the preference to %d", got)
			}
		})
	}
}

func TestPreferredColorMode(t *testing.T) {
	// The preferred color mode must be kept even on a platform without a window (#3480).
	defer ebiten.SetPreferredColorMode(ebiten.ColorModeUnknown)

	if got, want := ebiten.PreferredColorMode(), ebiten.ColorModeUnknown; got != want {
		t.Errorf("ebiten.PreferredColorMode(): got: %d, want: %d", got, want)
	}

	for _, mode := range []ebiten.ColorMode{ebiten.ColorModeDark, ebiten.ColorModeLight, ebiten.ColorModeUnknown} {
		ebiten.SetPreferredColorMode(mode)
		if got, want := ebiten.PreferredColorMode(), mode; got != want {
			t.Errorf("ebiten.PreferredColorMode() after ebiten.SetPreferredColorMode(%d): got: %d, want: %d", mode, got, want)
		}
	}
}
