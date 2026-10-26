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

package convert_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/audio/internal/convert"
)

type stallingReader struct {
	*bytes.Reader
	emptyReads int
	oneByte    bool
	calls      int
}

func (r *stallingReader) Read(p []byte) (int, error) {
	r.calls++
	if r.emptyReads > 0 {
		r.emptyReads--
		return 0, nil
	}
	if r.oneByte {
		p = p[:min(len(p), 1)]
	}
	n, err := r.Reader.Read(p)
	if r.oneByte && n > 0 {
		r.emptyReads = 99
	}
	return n, err
}

func TestConversionEmptyReads(t *testing.T) {
	factories := map[string]func(io.ReadSeeker) io.Reader{
		"float32":   func(r io.ReadSeeker) io.Reader { return convert.NewFloat32BytesReadSeekerFromInt16BytesReadSeeker(r) },
		"stereoF32": func(r io.ReadSeeker) io.Reader { return convert.NewStereoF32(r, false) },
		"monoF32":   func(r io.ReadSeeker) io.Reader { return convert.NewStereoF32(r, true) },
	}
	for _, mono := range []bool{false, true} {
		for _, format := range []convert.Format{convert.FormatU8, convert.FormatS16, convert.FormatS24} {
			factories[fmt.Sprintf("i16/mono=%v/format=%d", mono, format)] = func(r io.ReadSeeker) io.Reader {
				return convert.NewStereoI16ReadSeeker(r, mono, format)
			}
		}
	}
	for name, factory := range factories {
		t.Run(name, func(t *testing.T) {
			// Zero-valued PCM is finite in every supported source format.
			source := make([]byte, 24)
			want, err := io.ReadAll(factory(bytes.NewReader(source)))
			if err != nil {
				t.Fatal(err)
			}
			for _, emptyReads := range []int{0, 3, 99, 100} {
				t.Run(fmt.Sprint(emptyReads), func(t *testing.T) {
					reader := &stallingReader{Reader: bytes.NewReader(source), emptyReads: emptyReads}
					r := factory(reader)
					if n, err := r.Read(nil); n != 0 || err != nil || reader.calls != 0 {
						t.Fatal("empty destination read the source")
					}
					if emptyReads == 100 {
						if n, err := r.Read(make([]byte, 8)); n != 0 || !errors.Is(err, io.ErrNoProgress) || reader.calls != 100 {
							t.Errorf("stalled read: n=%d err=%v calls=%d", n, err, reader.calls)
						}
					}
					got, err := io.ReadAll(r)
					if err != nil || !bytes.Equal(got, want) {
						t.Errorf("resumed PCM: error=%v got %d bytes, want %d", err, len(got), len(want))
					}
				})
			}
			reader := &stallingReader{Reader: bytes.NewReader(source), oneByte: true}
			got, err := io.ReadAll(factory(reader))
			if err != nil || !bytes.Equal(got, want) {
				t.Errorf("partial progress did not reset empty-read counter: %v", err)
			}
		})
	}
}
