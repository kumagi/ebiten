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

package main

import (
	"math"
	"testing"
)

func TestObjectPointsEmpty(t *testing.T) {
	if points := (object{}).points(); len(points) != 0 {
		t.Errorf("empty object points = %v", points)
	}
	if rays := rayCasting(0, 0, []object{{}}); len(rays) != 0 {
		t.Errorf("empty object rays = %v", rays)
	}
}

func TestRayCastingOpenScene(t *testing.T) {
	for _, tc := range []struct {
		name    string
		objects []object
		count   int
	}{
		{"none", nil, 0},
		{"wall", []object{{walls: []line{{10, -10, 10, 10}}}}, 2},
		{"point", []object{{walls: []line{{10, 10, 10, 10}}}}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rays := rayCasting(0, 0, tc.objects)
			if len(rays) != tc.count {
				t.Errorf("rays = %v, want %d hits", rays, tc.count)
			}
			for _, ray := range rays {
				if math.Abs(ray.X2-10) > 1e-9 || ray.Y2 < -10 || ray.Y2 > 10 {
					t.Errorf("unexpected wall hit: %v", ray)
				}
			}
		})
	}
}

func TestRayCastingClosedScene(t *testing.T) {
	objects := []object{{walls: []line{{-10, -10, 10, -10}, {10, -10, 10, 10}, {10, 10, -10, 10}, {-10, 10, -10, -10}}}}
	rays := rayCasting(0, 0, objects)
	if len(rays) != 8 {
		t.Fatalf("rays = %d, want 8", len(rays))
	}
	for i, ray := range rays {
		if math.Abs(math.Max(math.Abs(ray.X2), math.Abs(ray.Y2))-10) > 1e-9 {
			t.Errorf("hit not on box: %v", ray)
		}
		if i > 0 && rays[i-1].angle() > ray.angle() {
			t.Errorf("rays not sorted at %d", i)
		}
	}
}
