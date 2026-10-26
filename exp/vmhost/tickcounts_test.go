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

package vmhost

import (
	"math"
	"testing"
)

func TestTickCounts(t *testing.T) {
	for _, count := range []int64{0, 1, 123, 1<<31 - 1, 1 << 31, 1<<32 + 7, 1 << 40} {
		want := count
		if want > int64(math.MaxInt) {
			want = int64(math.MaxInt)
		}
		for _, processed := range []bool{false, true} {
			g := &GuestSession{submittedTicks: count}
			if processed {
				g.consumedTicks = count
			}
			got := g.PendingTicks()
			if processed {
				got = g.ProcessedTicks()
			}
			if int64(got) != want {
				t.Errorf("count=%d processed=%v: got %d, want %d", count, processed, got, want)
			}
			if processed && g.PendingTicks() != 0 {
				t.Error("processed ticks still pending")
			}
		}
	}
}
