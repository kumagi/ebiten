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
	"encoding/json"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestKeyMarshalDefined(t *testing.T) {
	for key := ebiten.Key(0); key <= ebiten.KeyMax; key++ {
		b, err := key.MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != key.String() {
			t.Errorf("key %d: name=%q, want %q", key, b, key.String())
		}
		var got ebiten.Key
		if err := got.UnmarshalText(b); err != nil {
			t.Fatal(err)
		}
		if got != key {
			t.Errorf("key %d: decoded %d", key, got)
		}
	}
}

func TestKeyMarshalUndefined(t *testing.T) {
	maxInt := int(^uint(0) >> 1)
	for _, key := range []ebiten.Key{-1, ebiten.KeyMax + 1, ebiten.Key(maxInt), ebiten.Key(-maxInt - 1)} {
		if b, err := key.MarshalText(); err == nil || len(b) != 0 {
			t.Errorf("key %d: MarshalText=(%q,%v), want error", key, b, err)
		}
	}
	if b, err := json.Marshal(map[ebiten.Key]int{ebiten.KeyA: 1, -1: 2, ebiten.KeyMax + 1: 3}); err == nil || len(b) != 0 {
		t.Errorf("undefined JSON keys: output=%s error=%v", b, err)
	}
}
