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

package gl

import (
	"fmt"
	"strings"
	"syscall/js"
	"testing"
)

func TestTemporaryArrayInputLength(t *testing.T) {
	for name, call := range map[string]func(int, bool) js.Value{
		"uint8": func(minLength int, initialize bool) js.Value {
			var data []uint8
			if initialize {
				data = []uint8{3, 4}
			}
			return tmpUint8ArrayFromUint8Slice(minLength, data)
		},
		"float32": func(minLength int, initialize bool) js.Value {
			var data []float32
			if initialize {
				data = []float32{3, 4}
			}
			return tmpFloat32ArrayFromFloat32Slice(minLength, data)
		},
		"int32": func(minLength int, initialize bool) js.Value {
			var data []int32
			if initialize {
				data = []int32{3, 4}
			}
			return tmpInt32ArrayFromInt32Slice(minLength, data)
		},
	} {
		t.Run(name, func(t *testing.T) {
			for _, minLength := range []int{2, 32} {
				arr := call(minLength, true)
				if arr.Length() < minLength || arr.Index(0).Float() != 3 || arr.Index(1).Float() != 4 {
					t.Errorf("data not preserved at length %d", minLength)
				}
			}
			if arr := call(2, false); arr.Index(0).Float() != 3 || arr.Index(1).Float() != 4 {
				t.Error("nil input unexpectedly initialized the shared buffer")
			}
			before := tmpArrayBufferByteLength
			tmpUint8Array.SetIndex(0, 99)
			defer func() {
				if got := fmt.Sprint(recover()); !strings.Contains(got, "2 elements do not fit in a temporary buffer of 1") {
					t.Errorf("panic = %q, want input-length rejection", got)
				}
				if tmpArrayBufferByteLength != before || tmpUint8Array.Index(0).Int() != 99 {
					t.Error("rejected input modified the shared buffer")
				}
			}()
			call(1, true)
		})
	}
}
