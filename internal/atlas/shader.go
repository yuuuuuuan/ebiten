// Copyright 2020 The Ebiten Authors
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

package atlas

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2/internal/restorable"
	"github.com/hajimehoshi/ebiten/v2/internal/shaderir"
)

type Shader struct {
	ir     *shaderir.Program
	shader *restorable.Shader
}

func NewShader(ir *shaderir.Program) *Shader {
	// A shader is initialized lazily, and the lock is not needed.
	return &Shader{
		ir: ir,
	}
}

func (s *Shader) ensureShader() *restorable.Shader {
	if s.shader != nil {
		return s.shader
	}
	s.shader = restorable.NewShader(s.ir)
	s.ir = nil
	return s.shader
}

// MarkDisposed marks the shader as disposed. The actual operation is deferred.
// MarkDisposed can be called from finalizers.
//
// A function from finalizer must not be blocked, but disposing operation can be blocked.
// Defer this operation until it becomes safe. (#913)
func (s *Shader) MarkDisposed() {
	// As MarkDisposed can be invoked from finalizers, backendsM should not be used.
	deferredM.Lock()
	deferred = append(deferred, func() {
		s.dispose()
	})
	deferredM.Unlock()
}

func (s *Shader) dispose() {
	runtime.SetFinalizer(s, nil)
	if s.shader == nil {
		return
	}
	s.shader.Dispose()
	s.shader = nil
}

var (
	NearestFilterShader = &Shader{shader: restorable.NearestFilterShader}
	LinearFilterShader  = &Shader{shader: restorable.LinearFilterShader}
)
