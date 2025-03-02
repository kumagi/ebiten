// Copyright 2021 The Ebiten Authors
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

package audio

import (
	"io"
	"runtime"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio/internal/readerdriver"
)

type readerPlayerFactory struct {
	context    readerdriver.Context
	sampleRate int
}

var readerDriverForTesting readerdriver.Context

func newReaderPlayerFactory(sampleRate int) *readerPlayerFactory {
	f := &readerPlayerFactory{
		sampleRate: sampleRate,
	}
	if readerDriverForTesting != nil {
		f.context = readerDriverForTesting
	}
	// TODO: Consider the hooks.
	return f
}

type readerPlayer struct {
	context *Context
	player  readerdriver.Player
	src     io.Reader
	stream  *timeStream
	factory *readerPlayerFactory
	m       sync.Mutex
}

func (f *readerPlayerFactory) newPlayerImpl(context *Context, src io.Reader) (playerImpl, error) {
	p := &readerPlayer{
		src:     src,
		context: context,
		factory: f,
	}
	runtime.SetFinalizer(p, (*readerPlayer).Close)
	return p, nil
}

func (f *readerPlayerFactory) suspend() error {
	if f.context == nil {
		return nil
	}
	return f.context.Suspend()
}

func (f *readerPlayerFactory) resume() error {
	if f.context == nil {
		return nil
	}
	return f.context.Resume()
}

func (p *readerPlayer) ensurePlayer() error {
	// Initialize the underlying player lazily to enable calling NewContext in an 'init' function.
	// Accessing the underlying player functions requires the environment to be already initialized,
	// but if Ebiten is used for a shared library, the timing when init functions are called
	// is unexpectable.
	// e.g. a variable for JVM on Android might not be set.
	if p.factory.context == nil {
		c, ready, err := readerdriver.NewContext(p.factory.sampleRate, channelNum, bitDepthInBytes)
		if err != nil {
			return err
		}
		go func() {
			<-ready
			p.context.setReady()
		}()
		p.factory.context = c
	}
	if p.stream == nil {
		s, err := newTimeStream(p.src, p.factory.sampleRate)
		if err != nil {
			return err
		}
		p.stream = s
	}
	if p.player == nil {
		p.player = p.factory.context.NewPlayer(p.stream)
	}
	return nil
}

func (p *readerPlayer) Play() {
	p.m.Lock()
	defer p.m.Unlock()

	if err := p.ensurePlayer(); err != nil {
		p.context.setError(err)
		return
	}
	if p.player.IsPlaying() {
		return
	}
	p.player.Play()
	p.context.addPlayer(p)
}

func (p *readerPlayer) Pause() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.player == nil {
		return
	}
	if !p.player.IsPlaying() {
		return
	}

	p.player.Pause()
	p.context.removePlayer(p)
}

func (p *readerPlayer) IsPlaying() bool {
	p.m.Lock()
	defer p.m.Unlock()

	if p.player == nil {
		return false
	}
	return p.player.IsPlaying()
}

func (p *readerPlayer) Volume() float64 {
	p.m.Lock()
	defer p.m.Unlock()

	if err := p.ensurePlayer(); err != nil {
		p.context.setError(err)
		return 0
	}
	return p.player.Volume()
}

func (p *readerPlayer) SetVolume(volume float64) {
	p.m.Lock()
	defer p.m.Unlock()

	if err := p.ensurePlayer(); err != nil {
		p.context.setError(err)
		return
	}
	p.player.SetVolume(volume)
}

func (p *readerPlayer) Close() error {
	p.m.Lock()
	defer p.m.Unlock()
	runtime.SetFinalizer(p, nil)

	if p.player != nil {
		defer func() {
			p.player = nil
		}()
		p.player.Pause()
		return p.player.Close()
	}
	return nil
}

func (p *readerPlayer) Current() time.Duration {
	p.m.Lock()
	defer p.m.Unlock()
	if err := p.ensurePlayer(); err != nil {
		p.context.setError(err)
		return 0
	}

	sample := (p.stream.Current() - int64(p.player.UnplayedBufferSize())) / bytesPerSample
	return time.Duration(sample) * time.Second / time.Duration(p.factory.sampleRate)
}

func (p *readerPlayer) Rewind() error {
	return p.Seek(0)
}

func (p *readerPlayer) Seek(offset time.Duration) error {
	p.m.Lock()
	defer p.m.Unlock()

	if err := p.ensurePlayer(); err != nil {
		return err
	}

	if p.player.IsPlaying() {
		defer func() {
			p.player.Play()
		}()
	}
	p.player.Reset()
	return p.stream.Seek(offset)
}

func (p *readerPlayer) Err() error {
	p.m.Lock()
	defer p.m.Unlock()

	if p.player == nil {
		return nil
	}
	return p.player.Err()
}

func (p *readerPlayer) source() io.Reader {
	return p.src
}

type timeStream struct {
	r          io.Reader
	sampleRate int
	pos        int64

	// m is a mutex for this stream.
	// All the exported functions are protected by this mutex as Read can be read from a different goroutine than Seek.
	m sync.Mutex
}

func newTimeStream(r io.Reader, sampleRate int) (*timeStream, error) {
	s := &timeStream{
		r:          r,
		sampleRate: sampleRate,
	}
	if seeker, ok := s.r.(io.Seeker); ok {
		// Get the current position of the source.
		pos, err := seeker.Seek(0, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		s.pos = pos
	}
	return s, nil
}

func (s *timeStream) Read(buf []byte) (int, error) {
	s.m.Lock()
	defer s.m.Unlock()

	n, err := s.r.Read(buf)
	s.pos += int64(n)
	return n, err
}

func (s *timeStream) Seek(offset time.Duration) error {
	s.m.Lock()
	defer s.m.Unlock()

	o := int64(offset) * bytesPerSample * int64(s.sampleRate) / int64(time.Second)

	// Align the byte position with the samples.
	o -= o % bytesPerSample
	o += s.pos % bytesPerSample

	seeker, ok := s.r.(io.Seeker)
	if !ok {
		panic("audio: the source must be io.Seeker when seeking but not")
	}
	pos, err := seeker.Seek(o, io.SeekStart)
	if err != nil {
		return err
	}

	s.pos = pos
	return nil
}

func (s *timeStream) Current() int64 {
	s.m.Lock()
	defer s.m.Unlock()

	return s.pos
}
