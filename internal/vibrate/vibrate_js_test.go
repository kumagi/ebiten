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

package vibrate

import (
	"syscall/js"
	"testing"
	"time"
)

func TestVibrateNavigator(t *testing.T) {
	global := js.Global()
	object := global.Get("Object")
	original := object.Call("getOwnPropertyDescriptor", global, "navigator")
	defer func() {
		if original.IsUndefined() {
			global.Delete("navigator")
		} else {
			object.Call("defineProperty", global, "navigator", original)
		}
	}()

	for _, name := range []string{"undefined", "null", "no-vibrate", "vibrate"} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					t.Errorf("vibrate panicked: %v", p)
				}
			}()
			calls := 0
			navigator := object.New()
			callback := js.FuncOf(func(this js.Value, args []js.Value) any {
				calls++
				if !this.Equal(navigator) || len(args) != 1 || args[0].Float() != 12 {
					t.Errorf("unexpected vibration receiver or duration: %v, %v", this, args)
				}
				return true
			})
			defer callback.Release()
			switch name {
			case "undefined":
				navigator = js.Undefined()
			case "null":
				navigator = js.Null()
			case "vibrate":
				navigator.Set("vibrate", callback)
			}
			object.Call("defineProperty", global, "navigator", map[string]any{"configurable": true, "value": navigator})
			vibrate(12*time.Millisecond, 0.5)
			want := 0
			if name == "vibrate" {
				want = 1
			}
			if calls != want {
				t.Errorf("vibration calls = %d, want %d", calls, want)
			}
		})
	}
}
