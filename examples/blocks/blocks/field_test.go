// Copyright 2026 The Ebitengine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blocks_test

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2/examples/blocks/blocks"
)

func TestRotatePieceAtTop(t *testing.T) {
	field := &blocks.Field{}
	piece := blocks.Pieces[blocks.BlockType1]
	x, y := piece.InitialPosition()
	if got := field.RotatePieceRight(piece, x, y, blocks.Angle0); got != blocks.Angle0 {
		t.Errorf("RotatePieceRight at spawn: got %v, want %v", got, blocks.Angle0)
	}
	if got := field.RotatePieceLeft(piece, x, y, blocks.Angle0); got != blocks.Angle0 {
		t.Errorf("RotatePieceLeft at spawn: got %v, want %v", got, blocks.Angle0)
	}
	y = field.DropPiece(piece, x, y, blocks.Angle0)
	if got := field.RotatePieceRight(piece, x, y, blocks.Angle0); got != blocks.Angle90 {
		t.Errorf("RotatePieceRight after drop: got %v, want %v", got, blocks.Angle90)
	}
	if got := field.RotatePieceLeft(piece, x, y, blocks.Angle0); got != blocks.Angle270 {
		t.Errorf("RotatePieceLeft after drop: got %v, want %v", got, blocks.Angle270)
	}
}

func TestPiecesAtInitialPosition(t *testing.T) {
	field := &blocks.Field{}
	for typ, piece := range blocks.Pieces {
		x, y := piece.InitialPosition()
		if got := field.MovePieceToLeft(piece, x, y, blocks.Angle0); got != x-1 {
			t.Errorf("piece %v: MovePieceToLeft: got %v, want %v", typ, got, x-1)
		}
		if got := field.MovePieceToRight(piece, x, y, blocks.Angle0); got != x+1 {
			t.Errorf("piece %v: MovePieceToRight: got %v, want %v", typ, got, x+1)
		}
		if got := field.DropPiece(piece, x, y, blocks.Angle0); got != y+1 {
			t.Errorf("piece %v: DropPiece: got %v, want %v", typ, got, y+1)
		}
	}
}
