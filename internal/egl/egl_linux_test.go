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

package egl

import (
	"errors"
	"reflect"
	"testing"
)

func TestCloseReportsCleanupErrors(t *testing.T) {
	makeCurrent, destroyContext, destroySurface, terminate, getError := eglMakeCurrent, eglDestroyContext, eglDestroySurface, eglTerminate, eglGetError
	t.Cleanup(func() {
		eglMakeCurrent = makeCurrent
		eglDestroyContext = destroyContext
		eglDestroySurface = destroySurface
		eglTerminate = terminate
		eglGetError = getError
	})
	for mask := 0; mask < 16; mask++ {
		var calls []string
		var last int32
		reads := 0
		attempt := func(n int, name string) bool {
			calls = append(calls, name)
			if mask&(1<<n) != 0 {
				last = int32(0x3001 + n)
				return false
			}
			return true
		}
		eglMakeCurrent = func(d, draw, read, ctx uintptr) bool {
			if d != 1 || draw != 0 || read != 0 || ctx != 0 {
				t.Error("unexpected unbind arguments")
			}
			return attempt(0, "unbind")
		}
		eglDestroyContext = func(d, ctx uintptr) bool {
			if d != 1 || ctx != 3 {
				t.Error("unexpected context arguments")
			}
			return attempt(1, "context")
		}
		eglDestroySurface = func(d, surf uintptr) bool {
			if d != 1 || surf != 2 {
				t.Error("unexpected surface arguments")
			}
			return attempt(2, "surface")
		}
		eglTerminate = func(d uintptr) bool {
			if d != 1 {
				t.Error("unexpected display")
			}
			return attempt(3, "terminate")
		}
		eglGetError = func() int32 { reads++; code := last; last = Success; return code }
		c := &Context{display: 1, surface: 2, context: 3}
		err := c.Close()
		if !reflect.DeepEqual(calls, []string{"unbind", "context", "surface", "terminate"}) {
			t.Errorf("mask=%d calls=%v", mask, calls)
		}
		failures := 0
		for n := 0; n < 4; n++ {
			want := mask&(1<<n) != 0
			if want {
				failures++
			}
			if got := errors.Is(err, eglError(0x3001+n)); got != want {
				t.Errorf("mask=%d step=%d preserved=%v want=%v error=%v", mask, n, got, want, err)
			}
		}
		if reads != failures {
			t.Errorf("mask=%d error reads=%d want=%d", mask, reads, failures)
		}
		if c.display != 0 || c.surface != 0 || c.context != 0 {
			t.Error("cleanup handles retained")
		}
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
		if len(calls) != 4 {
			t.Error("second Close repeated cleanup")
		}
	}
}
func TestClosePartiallyCreated(t *testing.T) {
	makeCurrent, destroyContext, destroySurface, terminate := eglMakeCurrent, eglDestroyContext, eglDestroySurface, eglTerminate
	t.Cleanup(func() {
		eglMakeCurrent = makeCurrent
		eglDestroyContext = destroyContext
		eglDestroySurface = destroySurface
		eglTerminate = terminate
	})
	for _, tc := range []struct {
		display, context, surface uintptr
		want                      []string
	}{{0, 0, 0, nil}, {1, 0, 0, []string{"terminate"}}, {1, 3, 0, []string{"unbind", "context", "terminate"}}, {1, 0, 2, []string{"unbind", "surface", "terminate"}}} {
		var calls []string
		eglMakeCurrent = func(uintptr, uintptr, uintptr, uintptr) bool { calls = append(calls, "unbind"); return true }
		eglDestroyContext = func(uintptr, uintptr) bool { calls = append(calls, "context"); return true }
		eglDestroySurface = func(uintptr, uintptr) bool { calls = append(calls, "surface"); return true }
		eglTerminate = func(uintptr) bool { calls = append(calls, "terminate"); return true }
		c := &Context{display: tc.display, context: tc.context, surface: tc.surface}
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(calls, tc.want) {
			t.Errorf("got=%v want=%v", calls, tc.want)
		}
	}
}
