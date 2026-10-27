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

package imagebridge

import (
	"runtime"
	"testing"
)

func TestGetBeforeSet(t *testing.T) {
	saved := theBridge
	defer func() { theBridge = saved }()
	theBridge = nil
	defer func() {
		if got := recover(); got != "imagebridge: Get is called before Set" {
			t.Errorf("panic=%v, want Get-before-Set diagnostic", got)
		}
	}()
	Get[int]()
}

func TestGet(t *testing.T) {
	saved := theBridge
	defer func() { theBridge = saved }()
	Set(Bridge[int]{OriginalImage: func(i int) int { return i + 1 }})
	if got := Get[int]().OriginalImage(2); got != 3 {
		t.Errorf("bridge callback=%d,want 3", got)
	}
}

func TestGetWrongImageType(t *testing.T) {
	saved := theBridge
	defer func() { theBridge = saved }()
	Set(Bridge[int]{})
	defer func() {
		if got := recover(); got == nil {
			t.Error("Get with a different image type did not panic")
		} else if _, ok := got.(*runtime.TypeAssertionError); !ok {
			t.Errorf("panic type=%T,want *runtime.TypeAssertionError", got)
		}
	}()
	Get[string]()
}
