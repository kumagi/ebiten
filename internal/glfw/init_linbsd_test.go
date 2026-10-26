// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2026 The Ebitengine Authors

//go:build freebsd || linux || netbsd

package glfw_test

import (
	"slices"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
)

func TestParseUriList(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  []string
	}{
		{
			name: "empty",
			in:   "",
			out:  nil,
		},
		{
			name: "single file URI",
			in:   "file:///home/user/file.txt\r\n",
			out:  []string{"/home/user/file.txt"},
		},
		{
			name: "multiple lines",
			in:   "file:///a.txt\r\nfile:///b.txt\r\nfile:///c.txt\r\n",
			out:  []string{"/a.txt", "/b.txt", "/c.txt"},
		},
		{
			name: "comment lines skipped",
			in:   "#comment\r\nfile:///a.txt\r\n# another comment\r\n",
			out:  []string{"/a.txt"},
		},
		{
			name: "localhost authority",
			in:   "file://localhost/etc/fstab\r\n",
			out:  []string{"/etc/fstab"},
		},
		{
			name: "case insensitive localhost",
			in:   "file://LOCALHOST/a.txt\r\n",
			out:  []string{"/a.txt"},
		},
		{
			name: "unrecognized authority",
			in:   "file://other.example/a.txt\r\n",
		},
		{
			name: "mixed authorities",
			in:   "file:///a.txt\nfile://other.example/b.txt\nfile://localhost/c.txt\n",
			out:  []string{"/a.txt", "/c.txt"},
		},
		{
			name: "localhost with port",
			in:   "file://localhost:1234/a.txt\n",
		},
		{
			name: "percent escapes",
			in:   "file:///path%20with%20spaces/%E3%81%82.txt\r\n",
			out:  []string{"/path with spaces/あ.txt"},
		},
		{
			name: "invalid escapes kept as-is",
			in:   "/100%zz/x%2\r\n",
			out:  []string{"/100%zz/x%2"},
		},
		{
			name: "non-file URI passes through",
			in:   "http://example.com/index.html\r\n",
			out:  []string{"http://example.com/index.html"},
		},
		{
			name: "bare LF separators and trailing empty lines",
			in:   "file:///a.txt\nfile:///b.txt\n\n\n",
			out:  []string{"/a.txt", "/b.txt"},
		},
		{
			name: "file URI without a path dropped",
			in:   "file://localhost\r\nfile:///a.txt\r\n",
			out:  []string{"/a.txt"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := glfw.ParseUriList(tc.in); !slices.Equal(got, tc.out) {
				t.Errorf("parseUriList(%q) = %#v, want %#v", tc.in, got, tc.out)
			}
		})
	}
}
