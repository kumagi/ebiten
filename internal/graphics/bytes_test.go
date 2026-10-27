// Copyright 2023 The Ebitengine Authors
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

package graphics

import "testing"

func TestBytesPoolEmptyRequest(t *testing.T) {
	var pool bytesPool
	cached := make([]byte, 8)
	pool.put(cached)
	empty := pool.get(0)
	if empty.Len() != 0 || cap(empty.bytes) != 0 {
		t.Errorf("empty request: len=%d, cap=%d", empty.Len(), cap(empty.bytes))
	}
	empty.Release()
	if len(pool.pool) != 1 || pool.totalSize() != cap(cached) {
		t.Fatalf("empty request consumed cached buffer: count=%d, size=%d", len(pool.pool), pool.totalSize())
	}
	next := pool.get(4)
	if len(pool.pool) != 0 || next.Len() != 4 || &next.bytes[0] != &cached[0] {
		t.Error("nonempty request did not reuse the cached buffer")
	}
	next.Release()
	if len(pool.pool) != 1 || pool.totalSize() != cap(cached) {
		t.Error("released nonempty buffer was not returned to the pool")
	}
}
