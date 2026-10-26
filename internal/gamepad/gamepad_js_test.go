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

package gamepad_test

import (
	"encoding/hex"
	"syscall/js"
	"testing"
	"time"

	"github.com/hajimehoshi/ebiten/v2/internal/gamepad"
)

var object = js.Global().Get("Object")

// setNavigator replaces navigator with an object whose getGamepads calls f until the test ends.
func setNavigator(t *testing.T, f func() any) {
	orig := object.Call("getOwnPropertyDescriptor", js.Global(), "navigator")
	t.Cleanup(func() {
		if orig.IsUndefined() {
			js.Global().Delete("navigator")
			return
		}
		object.Call("defineProperty", js.Global(), "navigator", orig)
	})

	getGamepads := js.FuncOf(func(this js.Value, args []js.Value) any {
		return f()
	})
	t.Cleanup(getGamepads.Release)
	nav := object.New()
	nav.Set("getGamepads", getGamepads)
	desc := object.New()
	desc.Set("value", nav)
	desc.Set("configurable", true)
	object.Call("defineProperty", js.Global(), "navigator", desc)
}

// TestUpdatePausesPollingWithoutGamepads tests that update does not call navigator.getGamepads for one second after it finds no gamepad.
func TestUpdatePausesPollingWithoutGamepads(t *testing.T) {
	var calls int
	setNavigator(t, func() any {
		calls++
		return []any{nil, nil, nil, nil}
	})

	var gps gamepad.Gamepads
	n := &gamepad.JSGamepads{}
	if err := n.Init(&gps); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(n.RemoveEventListener)
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls after the first update: got %d, want 1", calls)
	}

	n.SetLastPoll(time.Now().Add(time.Hour))
	for range 10 {
		if err := n.Update(&gps); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 1 {
		t.Fatalf("calls during the pause: got %d, want 1", calls)
	}

	n.SetLastPoll(time.Now().Add(-time.Second))
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Errorf("calls after the pause: got %d, want 2", calls)
	}
}

// TestUpdateFindsGamepadAfterConnectedEvent tests that update finds a gamepad at the first update after a gamepadconnected event, even during the pause.
func TestUpdateFindsGamepadAfterConnectedEvent(t *testing.T) {
	var calls int
	var connected bool
	setNavigator(t, func() any {
		calls++
		if !connected {
			return []any{nil}
		}
		gp := object.New()
		gp.Set("index", 0)
		gp.Set("id", "test")
		gp.Set("mapping", "standard")
		gp.Set("axes", []any{})
		gp.Set("buttons", []any{})
		return []any{gp}
	})

	var gps gamepad.Gamepads
	n := &gamepad.JSGamepads{}
	if err := n.Init(&gps); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(n.RemoveEventListener)
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}

	n.SetLastPoll(time.Now().Add(time.Hour))
	connected = true
	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("gamepadconnected"))
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("calls after the event: got %d, want 2", calls)
	}
	if got := len(gps.AppendGamepadIDs(nil)); got != 1 {
		t.Errorf("gamepads after the event: got %d, want 1", got)
	}
}

// TestUpdateFindsGamepadWithoutEvent tests that update finds a gamepad that appears without a gamepadconnected event, then polls navigator.getGamepads at every update.
func TestUpdateFindsGamepadWithoutEvent(t *testing.T) {
	var calls int
	var connected, pressed bool
	setNavigator(t, func() any {
		calls++
		if !connected {
			return []any{nil}
		}
		button := object.New()
		button.Set("pressed", pressed)
		button.Set("value", 0)
		gp := object.New()
		gp.Set("index", 0)
		gp.Set("id", "test")
		gp.Set("mapping", "standard")
		gp.Set("axes", []any{})
		gp.Set("buttons", []any{button})
		return []any{gp}
	})

	var gps gamepad.Gamepads
	n := &gamepad.JSGamepads{}
	if err := n.Init(&gps); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(n.RemoveEventListener)
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}

	connected = true
	n.SetLastPoll(time.Now().Add(time.Hour))
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}
	if got := len(gps.AppendGamepadIDs(nil)); got != 0 {
		t.Fatalf("gamepads during the pause: got %d, want 0", got)
	}

	n.SetLastPoll(time.Now().Add(-time.Second))
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}
	ids := gps.AppendGamepadIDs(nil)
	if len(ids) != 1 {
		t.Fatalf("gamepads after the pause: got %d, want 1", len(ids))
	}
	if calls != 2 {
		t.Fatalf("calls after the pause: got %d, want 2", calls)
	}
	gp := gps.Get(ids[0])
	if gp.Button(0) {
		t.Fatalf("button 0 before the press: got true, want false")
	}

	n.SetLastPoll(time.Now().Add(time.Hour))
	pressed = true
	if err := n.Update(&gps); err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Errorf("calls after the press: got %d, want 3", calls)
	}
	if !gp.Button(0) {
		t.Errorf("button 0 after the press: got false, want true")
	}

	for range 10 {
		if err := n.Update(&gps); err != nil {
			t.Fatal(err)
		}
	}
	if calls != 13 {
		t.Errorf("calls after ten more updates: got %d, want 13", calls)
	}
}

func TestUpdateGamepadIndexReuse(t *testing.T) {
	for _, tc := range []struct {
		name    string
		id      string
		mapping string
	}{
		{
			name:    "id",
			id:      "second-controller",
			mapping: "",
		},
		{
			name:    "id suffix",
			id:      "first-controller-with-new-suffix",
			mapping: "",
		},
		{
			name:    "mapping",
			id:      "first-controller",
			mapping: "standard",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			id, mapping, pressed := "first-controller", "", false
			makeGamepad := func(index int, id, mapping string) js.Value {
				button := object.New()
				button.Set("pressed", pressed)
				button.Set("value", 0)
				gp := object.New()
				gp.Set("index", index)
				gp.Set("id", id)
				gp.Set("mapping", mapping)
				gp.Set("axes", []any{})
				gp.Set("buttons", []any{button})
				return gp
			}
			setNavigator(t, func() any { return []any{makeGamepad(0, id, mapping), nil, makeGamepad(2, "other-controller", "")} })
			var gps gamepad.Gamepads
			n := &gamepad.JSGamepads{}
			if err := n.Init(&gps); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(n.RemoveEventListener)
			if err := n.Update(&gps); err != nil {
				t.Fatal(err)
			}
			ids := gps.AppendGamepadIDs(nil)
			if len(ids) != 2 {
				t.Fatal("expected two gamepads")
			}
			first, other := gps.Get(ids[0]), gps.Get(ids[1])
			if err := n.Update(&gps); err != nil {
				t.Fatal(err)
			}
			if gps.Get(ids[0]) != first || gps.Get(ids[1]) != other {
				t.Error("unchanged metadata replaced a gamepad")
			}
			id, mapping, pressed = tc.id, tc.mapping, true
			if err := n.Update(&gps); err != nil {
				t.Fatal(err)
			}
			if got := len(gps.AppendGamepadIDs(nil)); got != 2 {
				t.Errorf("gamepad count=%d, want2", got)
			}
			current := gps.Get(ids[0])
			if current == first {
				t.Error("reused index retained the previous gamepad")
			}
			if got := current.Name(); got != tc.id {
				t.Errorf("name=%q, want%q", got, tc.id)
			}
			var guid [16]byte
			copy(guid[:], tc.id)
			if got, want := current.SDLID(), hex.EncodeToString(guid[:]); got != want {
				t.Errorf("SDLID=%q, want%q", got, want)
			}
			if got, want := current.IsStandardLayoutAvailable(), tc.mapping == "standard"; got != want {
				t.Errorf("standard layout=%v, want%v", got, want)
			}
			if !current.Button(0) {
				t.Error("replacement did not receive its button state")
			}
			if gps.Get(ids[1]) != other {
				t.Error("unrelated gamepad was replaced")
			}
		})
	}
}
