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

//go:build freebsd || linux || netbsd

package glfw_test

import (
	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
	"testing"
)

func TestConvertUTF8ToLatin1(t *testing.T) {
	var allLatin1 []rune
	var expected []byte
	for i := range 256 {
		allLatin1 = append(allLatin1, rune(i))
		expected = append(expected, byte(i))
	}
	for _, tc := range []struct{ input, want string }{
		{"", ""},
		{"ASCII\t\n", "ASCII\t\n"},
		{string(allLatin1), string(expected)},
		{"café 日本🚀", "caf\xe9 ???"},
		{"\xff", "?"},
	} {
		if got := glfw.ConvertUTF8ToLatin1(tc.input); got != tc.want {
			t.Errorf("conversion of %q: got %q, want %q", tc.input, got, tc.want)
		}
	}
}
