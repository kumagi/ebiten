// Copyright 2023 The Ebitengine Authors
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

package text

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type glyphBoundsCacheValue struct {
	bounds  fixed.Rectangle26_6
	advance fixed.Int26_6
	ok      bool
}

type glyphAdvanceCacheValue struct {
	advance fixed.Int26_6
	ok      bool
}

type kernCacheKey struct {
	r0 rune
	r1 rune
}

type faceWithCache struct {
	f       font.Face
	metrics font.Metrics

	// atime is the last time this face was used, as an ebiten tick. The
	// package-global caches hold an entry per face, so an abandoned face must
	// be collectable.
	atime int64

	glyphBoundsCache  map[rune]glyphBoundsCacheValue
	glyphAdvanceCache map[rune]glyphAdvanceCacheValue
	kernCache         map[kernCacheKey]fixed.Int26_6
}

func (f *faceWithCache) Close() error {
	if err := f.f.Close(); err != nil {
		return err
	}

	f.glyphBoundsCache = nil
	f.glyphAdvanceCache = nil
	f.kernCache = nil
	return nil
}

var faceWithCacheCache map[font.Face]*faceWithCache

func faceWithCacheFromFace(face font.Face) *faceWithCache {
	dropUnusedFaces()
	if f, ok := faceWithCacheCache[face]; ok {
		f.atime = ebiten.Tick()
		return f
	}

	f := &faceWithCache{
		f:     face,
		atime: ebiten.Tick(),
	}
	if faceWithCacheCache == nil {
		faceWithCacheCache = map[font.Face]*faceWithCache{}
	}
	faceWithCacheCache[face] = f
	return f
}

var lastFaceCacheCleanupTick int64 = -1

// dropUnusedFaces removes the faces that have not been used for a while, with
// their cached glyph images. The caches are package-global and keyed by the
// face, so a face the game stops using would otherwise keep its cached glyph
// images forever.
func dropUnusedFaces() {
	now := ebiten.Tick()
	if lastFaceCacheCleanupTick == now {
		return
	}
	lastFaceCacheCleanupTick = now
	for face, fc := range faceWithCacheCache {
		if fc.atime >= now-60 {
			continue
		}
		delete(glyphImageCache, fc)
		delete(faceWithCacheCache, face)
	}
}

func (f *faceWithCache) Glyph(dot fixed.Point26_6, r rune) (dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {
	// TODO: Move glyphImageCache to here.
	return f.f.Glyph(dot, r)
}

func (f *faceWithCache) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if v, ok := f.glyphBoundsCache[r]; ok {
		return v.bounds, v.advance, v.ok
	}

	bounds, advance, ok = f.f.GlyphBounds(r)
	if f.glyphBoundsCache == nil {
		f.glyphBoundsCache = map[rune]glyphBoundsCacheValue{}
	}
	f.glyphBoundsCache[r] = glyphBoundsCacheValue{
		bounds:  bounds,
		advance: advance,
		ok:      ok,
	}
	return
}

func (f *faceWithCache) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if v, ok := f.glyphAdvanceCache[r]; ok {
		return v.advance, v.ok
	}

	advance, ok = f.f.GlyphAdvance(r)
	if f.glyphAdvanceCache == nil {
		f.glyphAdvanceCache = map[rune]glyphAdvanceCacheValue{}
	}
	f.glyphAdvanceCache[r] = glyphAdvanceCacheValue{
		advance: advance,
		ok:      ok,
	}
	return
}

func (f *faceWithCache) Kern(r0, r1 rune) fixed.Int26_6 {
	key := kernCacheKey{r0: r0, r1: r1}
	if v, ok := f.kernCache[key]; ok {
		return v
	}

	v := f.f.Kern(r0, r1)
	if f.kernCache == nil {
		f.kernCache = map[kernCacheKey]fixed.Int26_6{}
	}
	f.kernCache[key] = v
	return v
}

func (f *faceWithCache) Metrics() font.Metrics {
	if f.metrics != (font.Metrics{}) {
		return f.metrics
	}
	f.metrics = f.f.Metrics()
	return f.metrics
}
