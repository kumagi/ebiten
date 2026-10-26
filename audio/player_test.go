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

package audio_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestPositionBeforeContextReady(t *testing.T) {
	for _, position := range []time.Duration{0, 500 * time.Millisecond} {
		t.Run(fmt.Sprint(position), func(t *testing.T) {
			setup()
			defer teardown()
			p, err := context.NewPlayerF32(bytes.NewReader(make([]byte, 44100*8)))
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := p.Close(); err != nil {
					t.Error(err)
				}
			}()
			if err := p.SetPosition(position); err != nil {
				t.Fatal(err)
			}
			p.Play()
			if err := context.UpdatePlayersForTesting(); err != nil {
				t.Fatal(err)
			}
			if got := p.Position(); got != position {
				t.Errorf("position before device creation: got %v want %v", got, position)
			}
			if err := context.InitDriverWithoutReadyForTesting(); err != nil {
				t.Fatal(err)
			}
			if context.IsReady() {
				t.Fatal("fixture delivered readiness early")
			}
			for range 3 {
				if err := context.UpdatePlayersForTesting(); err != nil {
					t.Fatal(err)
				}
				if got := p.Position(); got != position {
					t.Errorf("position before readiness: got %v want %v", got, position)
				}
			}
			context.MarkReadyForTesting()
			if !context.IsReady() {
				t.Error("readiness notification was not delivered")
			}
			if err := context.UpdatePlayersForTesting(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
