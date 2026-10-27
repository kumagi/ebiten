// Copyright 2014 Hajime Hoshi
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

package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/blocks/blocks"
)

var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	if *cpuProfile != "" {
		f, createErr := os.Create(*cpuProfile)
		if createErr != nil {
			return createErr
		}
		w := bufio.NewWriter(f)
		var profiling bool
		defer func() {
			if profiling {
				pprof.StopCPUProfile()
			}
			flushErr := w.Flush()
			closeErr := f.Close()
			if err == nil {
				err = flushErr
			}
			if err == nil {
				err = closeErr
			}
		}()
		if err := pprof.StartCPUProfile(w); err != nil {
			return err
		}
		profiling = true
	}

	ebiten.SetWindowSize(blocks.ScreenWidth*2, blocks.ScreenHeight*2)
	ebiten.SetWindowTitle("Blocks (Ebitengine Demo)")
	return ebiten.RunGame(&blocks.Game{})
}
