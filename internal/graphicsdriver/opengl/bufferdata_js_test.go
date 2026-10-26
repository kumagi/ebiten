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

//go:build !playstation5

package opengl_test

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl"
	"github.com/hajimehoshi/ebiten/v2/internal/graphicsdriver/opengl/gl"
	"strings"
	"testing"
)

const (
	invalidValue     = 0x0501
	invalidOperation = 0x0502
)

type bufferReadContext struct {
	gl.Context
	code          uint32
	reads, errors int
}

func (g *bufferReadContext) GetBufferSubData(target uint32, offset int, dst []byte) {
	if target != gl.PIXEL_PACK_BUFFER || offset != 2 {
		panic("unexpected read")
	}
	g.reads++
	copy(dst, []byte{1, 2, 3, 4})
}
func (g *bufferReadContext) GetError() uint32 { g.errors++; return g.code }
func TestReadBufferDataErrors(t *testing.T) {
	for _, code := range []uint32{gl.NO_ERROR, invalidOperation, invalidValue} {
		t.Run(fmt.Sprint(code), func(t *testing.T) {
			g := &bufferReadContext{code: code}
			dst := make([]byte, 4)
			err := opengl.ReadBufferDataForTesting(g, gl.PIXEL_PACK_BUFFER, 2, dst)
			if code == gl.NO_ERROR {
				if err != nil || !bytes.Equal(dst, []byte{1, 2, 3, 4}) {
					t.Errorf("normal read %v: %v", dst, err)
				}
			} else if err == nil || !strings.Contains(err.Error(), fmt.Sprint(code)) {
				t.Errorf("error %v does not report code %d", err, code)
			}
			if g.reads != 1 || g.errors != 1 {
				t.Errorf("reads=%d error checks=%d", g.reads, g.errors)
			}
		})
	}
	g := &bufferReadContext{code: invalidOperation}
	if err := opengl.ReadBufferDataForTesting(g, gl.PIXEL_PACK_BUFFER, 2, nil); err != nil || g.reads != 0 || g.errors != 0 {
		t.Error("empty read touched GL state")
	}
}
