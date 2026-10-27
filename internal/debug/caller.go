// Copyright 2024 The Ebitengine Authors
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

package debug

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Use this file's recorded path so the root matches runtime frames even with -trimpath.
var ebitengineFileDir = func() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return path.Dir(path.Dir(path.Dir(filepath.ToSlash(file))))
}()

type CallerType int

const (
	CallerTypeNone CallerType = iota
	CallerTypeRegular
	CallerTypeInternal
)

// FirstCaller returns the file and line number of the first caller outside of Ebitengine.
func FirstCaller() (file string, line int, callerType CallerType) {
	if ebitengineFileDir == "" {
		return "", 0, CallerTypeNone
	}

	// Relying on a caller stacktrace is very fragile, but this is fine as this is only for debugging.
	var ebitenPackageReached bool
	for i := 0; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// The file should be with a slash, but just in case, convert it.
		file = filepath.ToSlash(file)

		if !ebitenPackageReached {
			if path.Dir(file) == ebitengineFileDir {
				ebitenPackageReached = true
			}
			continue
		}

		if path.Dir(file) == ebitengineFileDir {
			continue
		}
		if strings.HasPrefix(path.Dir(file), ebitengineFileDir+"/") && !strings.HasPrefix(path.Dir(file), ebitengineFileDir+"/examples/") {
			if strings.HasPrefix(path.Dir(file), ebitengineFileDir+"/internal/") {
				return file, line, CallerTypeInternal
			}
			continue
		}
		return file, line, CallerTypeRegular
	}

	return "", 0, CallerTypeNone
}
