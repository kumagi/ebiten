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
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

func TestKeyStringModifiers(t *testing.T) {
	for _, tc := range []struct {
		key  ui.Key
		want string
	}{
		{ui.KeyAlt, "KeyAlt"}, {ui.KeyControl, "KeyControl"},
		{ui.KeyShift, "KeyShift"}, {ui.KeyMeta, "KeyMeta"}, {ui.KeyMax, "KeyMeta"},
		{ui.KeyA, "KeyA"}, {ui.KeyAltLeft, "KeyAltLeft"},
		{-1, "Key(-1)"}, {ui.KeyMax + 1, fmt.Sprintf("Key(%d)", ui.KeyMax+1)},
	} {
		if got := tc.key.String(); got != tc.want {
			t.Errorf("key %d: got %q, want %q", tc.key, got, tc.want)
		}
	}
}
