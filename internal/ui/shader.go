// Copyright 2022 The Ebiten Authors
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

package ui

import (
	"fmt"
	"math"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2/internal/atlas"
	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/shaderir"
)

type Shader struct {
	shader *atlas.Shader

	uniformNames      []string
	uniformTypes      []shaderir.Type
	uniformDwordCount int
}

func NewShader(ir *shaderir.Program, name string) *Shader {
	uniformTypes := ir.Uniforms[graphics.PreservedUniformVariablesCount:]
	var uniformDwordCount int
	for _, typ := range uniformTypes {
		uniformDwordCount += typ.DwordCount()
	}
	return &Shader{
		shader:            atlas.NewShader(ir, name),
		uniformNames:      ir.UniformNames[graphics.PreservedUniformVariablesCount:],
		uniformTypes:      uniformTypes,
		uniformDwordCount: uniformDwordCount,
	}
}

func (s *Shader) Deallocate() {
	s.shader.Deallocate()
}

func uniformIsSequence(typ shaderir.Type) bool {
	switch typ.Main {
	case shaderir.Vec2, shaderir.Vec3, shaderir.Vec4,
		shaderir.IVec2, shaderir.IVec3, shaderir.IVec4,
		shaderir.Mat2, shaderir.Mat3, shaderir.Mat4, shaderir.Array:
		return true
	default:
		return false
	}
}

func uniformScalarType(typ shaderir.Type) shaderir.Type {
	for typ.Main == shaderir.Array {
		if len(typ.Sub) == 0 {
			return shaderir.Type{}
		}
		typ = typ.Sub[0]
	}
	switch typ.Main {
	case shaderir.Vec2, shaderir.Vec3, shaderir.Vec4,
		shaderir.Mat2, shaderir.Mat3, shaderir.Mat4:
		return shaderir.Type{Main: shaderir.Float}
	case shaderir.IVec2, shaderir.IVec3, shaderir.IVec4:
		return shaderir.Type{Main: shaderir.Int}
	default:
		return typ
	}
}

func uniformValueMatches(typ shaderir.Type, v reflect.Value) bool {
	switch typ.Main {
	case shaderir.Bool:
		return v.Kind() == reflect.Bool
	case shaderir.Int:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64:
			return true
		}
	case shaderir.Float:
		switch v.Kind() {
		case reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return true
		}
	}
	return false
}

func panicUniform(name string, typ shaderir.Type) {
	panic(fmt.Sprintf("ui: unexpected uniform value for %s (%s)", name, typ.String()))
}

func validateUniformValue(name string, typ shaderir.Type, v reflect.Value) {
	scalarType := uniformScalarType(typ)
	if (scalarType.Main != shaderir.Bool && scalarType.Main != shaderir.Int && scalarType.Main != shaderir.Float) || !v.IsValid() {
		panicUniform(name, typ)
	}
	if uniformIsSequence(typ) {
		if (v.Kind() != reflect.Slice && v.Kind() != reflect.Array) || v.Len() != typ.DwordCount() {
			panicUniform(name, typ)
		}
		for i := range v.Len() {
			if !uniformValueMatches(scalarType, v.Index(i)) {
				panicUniform(name, typ)
			}
		}
		return
	}
	if (v.Kind() == reflect.Slice || v.Kind() == reflect.Array) || !uniformValueMatches(scalarType, v) {
		panicUniform(name, typ)
	}
}

func uniformFloatValue(v reflect.Value) float64 {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	default:
		return float64(v.Uint())
	}
}

func writeUniformValue(dst []uint32, idx int, typ shaderir.Type, v reflect.Value) {
	if uniformIsSequence(typ) {
		elementType := uniformScalarType(typ)
		for i := range v.Len() {
			writeUniformValue(dst, idx+i, elementType, v.Index(i))
		}
		return
	}
	switch typ.Main {
	case shaderir.Bool:
		if v.Bool() {
			dst[idx] = 1
		} else {
			dst[idx] = 0
		}
	case shaderir.Int, shaderir.IVec2, shaderir.IVec3, shaderir.IVec4:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dst[idx] = uint32(v.Int())
		case reflect.Float32, reflect.Float64:
			dst[idx] = uint32(int64(v.Float()))
		default:
			dst[idx] = uint32(v.Uint())
		}
	case shaderir.Float, shaderir.Vec2, shaderir.Vec3, shaderir.Vec4,
		shaderir.Mat2, shaderir.Mat3, shaderir.Mat4:
		dst[idx] = math.Float32bits(float32(uniformFloatValue(v)))
	default:
		panic("ui: unexpected uniform type")
	}
}

func (s *Shader) AppendUniforms(dst []uint32, uniforms map[string]any) []uint32 {
	origLen := len(dst)
	if cap(dst)-len(dst) >= s.uniformDwordCount {
		dst = dst[:len(dst)+s.uniformDwordCount]
		for i := origLen; i < len(dst); i++ {
			dst[i] = 0
		}
	} else {
		dst = append(dst, make([]uint32, s.uniformDwordCount)...)
	}

	idx := origLen
	for i, name := range s.uniformNames {
		typ := s.uniformTypes[i]
		// Ignore unused names (#2710).
		if uv, ok := uniforms[name]; ok {
			v := reflect.ValueOf(uv)
			validateUniformValue(name, typ, v)
			writeUniformValue(dst, idx, typ, v)
		}
		idx += typ.DwordCount()
	}

	return dst
}
