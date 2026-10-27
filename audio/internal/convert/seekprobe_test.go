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

type failingEndProbeSource struct {
	*bytes.Reader
	probeErr, restoreErr error
}

func (s *failingEndProbeSource) Seek(offset int64, whence int) (int64, error) {
	if whence == io.SeekEnd {
		n, _ := s.Reader.Seek(offset, whence)
		return n, s.probeErr
	}
	if whence == io.SeekStart && s.restoreErr != nil {
		return 0, s.restoreErr
	}
	return s.Reader.Seek(offset, whence)
}

func TestStereoSeekEndRestoresFailedProbe(t *testing.T) {
	constructors := map[string]func(io.ReadSeeker) io.ReadSeeker{}
	for _, mono := range []bool{false, true} {
		constructors[fmt.Sprintf("F32/mono=%v", mono)] = func(src io.ReadSeeker) io.ReadSeeker { return convert.NewStereoF32(src, mono) }
		for _, format := range []convert.Format{convert.FormatU8, convert.FormatS16, convert.FormatS24} {
			constructors[fmt.Sprintf("I16/%v/mono=%v", format, mono)] = func(src io.ReadSeeker) io.ReadSeeker { return convert.NewStereoI16ReadSeeker(src, mono, format) }
		}
	}
	data := make([]byte, 96)
	for i := range data {
		data[i] = byte(i)
	}
	probeErr, restoreErr := errors.New("end probe"), errors.New("restore")
	for name, newReader := range constructors {
		for _, failRestore := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/restore-fails=%v", name, failRestore), func(t *testing.T) {
				src := &failingEndProbeSource{Reader: bytes.NewReader(data), probeErr: probeErr}
				if failRestore {
					src.restoreErr = restoreErr
				}
				r := newReader(src)
				// Leave a complete source frame buffered without consuming any output.
				if _, err := r.Read(make([]byte, 1)); err != io.ErrShortBuffer {
					t.Fatal(err)
				}
				before, _ := src.Reader.Seek(0, io.SeekCurrent)
				_, err := r.Seek(0, io.SeekEnd)
				if !errors.Is(err, probeErr) {
					t.Errorf("lost probe error: %v", err)
				}
				if failRestore {
					if !errors.Is(err, restoreErr) {
						t.Errorf("lost restore error: %v", err)
					}
					return
				}
				if after, _ := src.Reader.Seek(0, io.SeekCurrent); after != before {
					t.Errorf("source position=%d want%d", after, before)
				}
				got, err := io.ReadAll(r)
				if err != nil {
					t.Fatal(err)
				}
				control := newReader(bytes.NewReader(data))
				if _, err := control.Read(make([]byte, 1)); err != io.ErrShortBuffer {
					t.Fatal(err)
				}
				want, err := io.ReadAll(control)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want) {
					t.Errorf("buffered stream changed after failed probe: got %d bytes want%d", len(got), len(want))
				}
			})
		}
	}
}
