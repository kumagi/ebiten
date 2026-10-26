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
	"errors"
	"os"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/glfw"
)

func TestX11Smoke(t *testing.T) {
	if os.Getenv("DISPLAY") == "" {
		t.Skip("no DISPLAY is available")
	}
	if err := glfw.SmokeTestX11(); err != nil {
		t.Fatal(err)
	}
}

func TestMonitorPositionFailureAndFallback(t *testing.T) {
	for _, tc := range []struct {
		name                         string
		available, broken, wantError bool
	}{
		{name: "RandR unavailable"},
		{name: "RandR broken", available: true, broken: true},
		{name: "resources unavailable", available: true, wantError: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			x, y, err := glfw.MonitorPositionForTesting(tc.available, tc.broken)
			if tc.wantError {
				if !errors.Is(err, glfw.PlatformError) {
					t.Errorf("GetPos error = %v, want PlatformError", err)
				}
			} else if x != 0 || y != 0 || err != nil {
				t.Errorf("fallback GetPos = (%d, %d, %v), want (0, 0, nil)", x, y, err)
			}
		})
	}
}
