package main

import (
	"fmt"
	"os"

	"github.com/acsellers/glex"
	"github.com/acsellers/math3d"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

const (
	vertex = `#version 130
#extension GL_ARB_explicit_attrib_location : enable
layout(location = 0) in vec3 vertexPosition_modelspace;

uniform mat4 MVP;

void main() {
  gl_Position = MVP * vec4(vertexPosition_modelspace, 1);
}`
	fragment = `#version 130
out vec3 color;

void main() {
  color = vec3(1,0,0);
}`
)

var (
	Window glex.ContextImplementation
)

type TriActivity struct {
	ctx *glex.Context
	buf gl.Buffer
	mvp *math3d.Matrix
}

func (ta *TriActivity) SetContext(c *glex.Context) {
	ta.ctx = c
}

func (ta *TriActivity) Draw() {
	ta.ctx.Program.Use()

	al := gl.AttribLocation(0)
	al.EnableArray()
	ta.buf.Bind(gl.ARRAY_BUFFER)
	al.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil,
	)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	al.DisableArray()
}

func (ta *TriActivity) Camera() {
	matrix := ta.ctx.Program.GetUniformLocation("MVP")
	matrix.UniformMatrix4fv(false, ta.mvp.Values())
}

func (ta *TriActivity) Generate() {
	ta.buf = gl.GenBuffer()
	ta.buf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		9*4,
		[]float32{
			-1.0, -1.0, 0.0,
			1.0, -1.0, 0.0,
			0.0, 1.0, 0.0,
		},
		gl.STATIC_DRAW,
	)

	proj := math3d.Perspective(45.0, 4.0/3.0, 0.1, 100.0)
	view := math3d.LookAt(
		math3d.NewVector([3]float32{4.0, 3.0, 3.0}),
		math3d.NewVector([3]float32{0, 0, 0}),
		math3d.NewVector([3]float32{0, 1.0, 0}),
	)
	model := math3d.Identity()

	ta.mvp = proj.Multiply(view.Multiply(model))
}

func init() {
	Window = &glex.GlfwContext{
		Ctx: &glex.Context{
			Width:       800,
			Height:      600,
			ColorDepth:  [4]int{8, 8, 8, 8},
			BufferDepth: [2]int{0, 0},
			Vsync:       true,
			WindowTitle: "Experiment: Four Refactor",
			Keyboard: func(key, state int) {
				if key == glfw.KeyEsc {
					Window.Context().Active = false
				}
			},
			Shaders: [2]string{
				vertex,
				fragment,
			},
		},
		Act: &TriActivity{},
	}
}

func main() {
	err := Window.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gl.GenVertexArray().Bind()

	for Window.Active() {
		Window.Swap()
	}
}
