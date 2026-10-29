// Copyright 2025 The Ebitengine Authors
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

package file_test

import (
	"io"
	"io/fs"
	"syscall/js"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/file"
)

func TestReadDirPartialError(t *testing.T) {
	for _, count := range []int{-1, 2} {
		var callbacks []js.Func
		callback := func(f func(js.Value, []js.Value) any) js.Func {
			cb := js.FuncOf(f)
			callbacks = append(callbacks, cb)
			return cb
		}
		var failed bool
		makeEntry := func(name string, failOnce bool) js.Value {
			return js.ValueOf(map[string]any{
				"name": name, "isFile": true, "isDirectory": false,
				"file": callback(func(_ js.Value, args []js.Value) any {
					if failOnce && !failed {
						failed = true
						args[1].Invoke(js.Global().Get("Error").New("file is temporarily unavailable"))
					} else {
						args[0].Invoke(js.ValueOf(map[string]any{"size": 1, "lastModified": 0}))
					}
					return nil
				}),
			})
		}
		entries := []any{makeEntry("first", false), makeEntry("second", true)}
		var listed bool
		reader := js.ValueOf(map[string]any{"readEntries": callback(func(_ js.Value, args []js.Value) any {
			if listed {
				args[0].Invoke(js.ValueOf([]any{}))
			} else {
				listed = true
				args[0].Invoke(js.ValueOf(entries))
			}
			return nil
		})})
		root := js.ValueOf(map[string]any{"name": "drop", "fullPath": "/drop", "createReader": callback(func(_ js.Value, _ []js.Value) any { return reader })})
		vfs, err := file.NewFileEntryFS([]js.Value{root})
		if err != nil {
			t.Fatal(err)
		}
		opened, err := vfs.Open(".")
		if err != nil {
			t.Fatal(err)
		}
		dir := opened.(fs.ReadDirFile)
		got, err := dir.ReadDir(count)
		if err == nil {
			t.Errorf("ReadDir(%d): missing error", count)
		}
		if len(got) != 1 || got[0].Name() != "first" {
			t.Errorf("ReadDir(%d): got %v, want first entry with error", count, got)
		}
		got, err = dir.ReadDir(count)
		if err != nil || len(got) != 1 || got[0].Name() != "second" {
			t.Errorf("ReadDir(%d) retry: got %v, %v, want only second entry", count, got, err)
		}
		if got, err := dir.ReadDir(1); len(got) != 0 || err != io.EOF {
			t.Errorf("end: got %v, %v, want EOF", got, err)
		}
		_ = opened.Close()
		for _, cb := range callbacks {
			cb.Release()
		}
	}
}
