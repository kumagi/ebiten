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

package hlsl_test

import (
	"strings"
	"testing"

	"github.com/hajimehoshi/ebiten/v2/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2/internal/shader"
	"github.com/hajimehoshi/ebiten/v2/internal/shaderir/hlsl"
)

func TestVertexAttributeTypes(t *testing.T) {
	p, err := shader.Compile([]byte(`package main

func Vertex(position vec4, texCoord vec3, alpha float, extra vec2) vec4 {
 return vec4(position.xy, position.z + texCoord.z + alpha + extra.x, position.w)
}

func Fragment(position vec4) vec4 {
 return vec4(1)
}
`), "Vertex", "Fragment", 0)
	if err != nil {
		t.Fatal(err)
	}
	vs, _, _, _ := hlsl.Compile(p)
	want := "Varyings VSMain(float4 A0 : POSITION, float3 A1 : TEXCOORD, float A2 : COLOR0, float2 A3 : COLOR1)"
	if !strings.Contains(vs, want) {
		t.Errorf("vertex signature does not preserve IR types; want %q in:\n%s", want, vs)
	}
}

func TestStandardVertexAttributeTypes(t *testing.T) {
	p, err := graphics.CompileShader([]byte(`//kage:unit pixels

package main

func Fragment(position vec4, srcPos vec2, color vec4) vec4 {
 return color
}
`))
	if err != nil {
		t.Fatal(err)
	}
	vs, _, _, _ := hlsl.Compile(p)
	want := "Varyings VSMain(float2 A0 : POSITION, float2 A1 : TEXCOORD, float4 A2 : COLOR0, float4 A3 : COLOR1)"
	if !strings.Contains(vs, want) {
		t.Errorf("standard vertex signature changed; want %q in:\n%s", want, vs)
	}
}
