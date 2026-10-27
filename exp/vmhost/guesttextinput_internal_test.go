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
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/vmprotocol"
)

func TestEndWithErrorNilPreservesSession(t *testing.T) {
	pooled := &vmprotocol.HostMessage{}
	g := &GuestSession{msgPool: []*vmprotocol.HostMessage{pooled}}
	g.cond = sync.NewCond(&g.mu)
	ti := &GuestTextInput{g: g, id: 42}
	func() {
		defer func() {
			if got, want := fmt.Sprint(recover()), "vmhost: EndWithError requires a non-nil error"; got != want {
				t.Errorf("panic = %q, want %q", got, want)
			}
		}()
		ti.EndWithError(nil)
	}()
	if ti.IsClosed() || len(g.ops) != 0 || len(g.msgPool) != 1 || g.msgPool[0] != pooled || pooled.Kind != 0 {
		t.Error("nil error changed session state or consumed a message")
	}
}

func TestEndWithErrorQueuesFailure(t *testing.T) {
	g := &GuestSession{}
	g.cond = sync.NewCond(&g.mu)
	ti := &GuestTextInput{g: g, id: 42}
	ti.EndWithError(errors.New("IME failure"))
	if !ti.IsClosed() || len(g.ops) != 1 {
		t.Fatalf("closed = %v, queued = %d", ti.IsClosed(), len(g.ops))
	}
	msg := g.ops[0].msg
	if g.ops[0].kind != opMessage || msg.Kind != vmprotocol.HostMessageKindTextInputState || msg.TextInputID != 42 || msg.TextInputState.Err != "IME failure" {
		t.Errorf("unexpected error message: %+v", msg)
	}
}
