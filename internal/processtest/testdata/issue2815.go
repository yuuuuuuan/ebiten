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

//go:build ignore

package main

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	init  bool
	count int
	end0  chan struct{}
	end1  chan struct{}
}

func (g *Game) Update() error {
	if !g.init {
		g.end0 = make(chan struct{})
		g.end1 = make(chan struct{})
		img := ebiten.NewImage(1, 1)
		go func() {
			t := time.Tick(time.Microsecond)
		loop:
			for {
				select {
				case <-t:
					img.At(0, 0)
				case <-g.end0:
					close(g.end1)
					break loop
				}
			}
		}()
		g.init = true
	}
	g.count++
	if g.count >= 60 {
		close(g.end0)
		<-g.end1
		return ebiten.Termination
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(w, h int) (int, int) {
	return 320, 240
}

func main() {
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
