// Copyright 2024 The Ebitengine Authors
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

package main_test

import (
	"path/filepath"
	"testing"

	ebitenmobile "github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile"
)

func TestJavaPackageName(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  "",
			out: false,
		},
		{
			in:  ".",
			out: false,
		},
		{
			in:  "com.hajimehoshi.goinovation",
			out: true,
		},
		{
			in:  "com.hajimehoshi.$goinovation",
			out: true,
		},
		{
			in:  "com.hajimehoshi..goinovation",
			out: false,
		},
		{
			in:  "com.hajimehoshi.go-inovation",
			out: false,
		},
		{
			in:  "com.hajimehoshi.strictfp", // strictfp is a Java keyword.
			out: false,
		},
		{
			in:  "com.hajimehoshi.null",
			out: false,
		},
		{
			in:  "com.hajimehoshi.go1inovation",
			out: true,
		},
		{
			in:  "com.hajimehoshi.1goinovation",
			out: false,
		},
		{
			in:  "あ.いうえお",
			out: true,
		},
	}
	for _, tc := range testCases {
		if got, want := ebitenmobile.IsValidJavaPackageName(tc.in), tc.out; got != want {
			t.Errorf("isValidJavaPackageName(%q) = %v; want %v", tc.in, got, want)
		}
	}
}

func TestIOSVersion(t *testing.T) {
	testCases := []struct {
		in  string
		out bool
	}{
		{
			in:  "12",
			out: false,
		},
		{
			in:  "12.99",
			out: false,
		},
		{
			in:  "13",
			out: true,
		},
		{
			in:  "13.0",
			out: true,
		},
		{
			in:  "13.0.1",
			out: true,
		},
		{
			in:  "14.1",
			out: true,
		},
		{
			in:  "26.0",
			out: true,
		},
		{
			in:  "",
			out: false,
		},
		{
			in:  "invalid",
			out: false,
		},
		{
			in:  "+Inf",
			out: false,
		},
		{
			in:  "NaN",
			out: false,
		},
		{
			in:  "+13",
			out: false,
		},
		{
			in:  "13e1",
			out: false,
		},
		{
			in:  "13.",
			out: false,
		},
		{
			in:  "13..0",
			out: false,
		},
		{
			in:  "13.0.0.1",
			out: false,
		},
		{
			in:  "13.0 ",
			out: false,
		},
		{
			in:  "13.999999999999999999999999",
			out: false,
		},
	}
	for _, tc := range testCases {
		if got := ebitenmobile.IsValidIOSVersion(tc.in); got != tc.out {
			t.Errorf("IsValidIOSVersion(%q) = %v; want %v", tc.in, got, tc.out)
		}
	}
}

func TestFrameworkName(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{
			in:  "MyGame.xcframework",
			out: "MyGame",
		},
		{
			in:  "mygame.xcframework",
			out: "Mygame",
		},
		{
			in:  "MG.xcframework",
			out: "MG",
		},
		{
			in:  "my_game.xcframework",
			out: "My_game",
		},
		{
			in:  filepath.Join("path", "to", "MyGame.xcframework"),
			out: "MyGame",
		},
	}
	for _, tc := range testCases {
		if got, want := ebitenmobile.FrameworkName(tc.in), tc.out; got != want {
			t.Errorf("frameworkName(%q) = %q; want %q", tc.in, got, want)
		}
	}
}
