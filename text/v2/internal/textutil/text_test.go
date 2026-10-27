// Copyright 2025 The Ebitengine Authors
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

package textutil_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/text/v2/internal/textutil"
)

func TestLines(t *testing.T) {
	testCases := []struct {
		In  string
		Out []string
	}{
		{
			In:  "",
			Out: nil,
		},
		{
			In:  "\n",
			Out: []string{"\n"},
		},
		{
			In:  "aaa\nbbb\nccc",
			Out: []string{"aaa\n", "bbb\n", "ccc"},
		},
		{
			In:  "aaa\nbbb\nccc\n",
			Out: []string{"aaa\n", "bbb\n", "ccc\n"},
		},
		{
			In:  "aaa\u0085bbb\r\nccc\rddd\u2028eee",
			Out: []string{"aaa\u0085", "bbb\r\n", "ccc\r", "ddd\u2028", "eee"},
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.In), func(t *testing.T) {
			lines := textutil.Lines(tc.In)
			for iteration := range 2 {
				got := slices.Collect(lines)
				if !slices.Equal(got, tc.Out) {
					t.Errorf("iteration %d: got %q, want %q", iteration, got, tc.Out)
				}
			}
		})
	}
}

func TestTrimTailingLineBreak(t *testing.T) {
	testCases := []struct {
		In  string
		Out string
	}{
		{
			In:  "",
			Out: "",
		},
		{
			In:  "aaa",
			Out: "aaa",
		},
		{
			In:  "aaa\n",
			Out: "aaa",
		},
		{
			In:  "aaa\n\n",
			Out: "aaa\n",
		},
		{
			In:  "aaa\r\n",
			Out: "aaa",
		},
		{
			In:  "aaa\r",
			Out: "aaa",
		},
		{
			In:  "aaa\u0085",
			Out: "aaa",
		},
		{
			In:  "aaa\u0085\u2028",
			Out: "aaa\u0085",
		},
		{
			In:  "aaa\nbbb\n",
			Out: "aaa\nbbb",
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%q", tc.In), func(t *testing.T) {
			got := textutil.TrimTailingLineBreak(tc.In)
			want := tc.Out
			if got != want {
				t.Errorf("got: %q, want: %q", got, want)
			}
		})
	}
}

func TestLinesAfterEarlyStop(t *testing.T) {
	lines := textutil.Lines("a\nb\nc\nd")
	var count int
	for range lines {
		count++
		if count == 2 {
			break
		}
	}
	got := slices.Collect(lines)
	want := []string{"a\n", "b\n", "c\n", "d"}
	if !slices.Equal(got, want) {
		t.Errorf("after early stop: got %q, want %q", got, want)
	}
}
