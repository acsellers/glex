package main

import (
	"fmt"
	"os"

	"github.com/acsellers/glex"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

var (
	Window  glex.ContextImplementation
	Ctx     *glex.Context
	Current glex.Activity
)

type TriActivity struct {
	ctx  *glex.Context
	buf  gl.Buffer
	gen  bool
	prog gl.Program
}

func (ta *TriActivity) SetContext(c *glex.Context) {
	ta.ctx = c
}

func (ta *TriActivity) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	ta.prog.Use()
	ta.DrawTriangle()
}

func (ta *TriActivity) DrawTriangle() {
	if !ta.gen {
		ta.GenBuffer()
		ta.GenShaders()
	}

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

func (ta *TriActivity) GenBuffer() {
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
	ta.gen = true
}

func (ta *TriActivity) GenShaders() {
	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	defer vShader.Delete()
	vShader.Source(`#version 130
#extension GL_ARB_explicit_attrib_location : enable
layout(location = 0) in vec3 vertexPosition_modelspace;

void main() {
  gl_Position.xyz = vertexPosition_modelspace;
  gl_Position.w = 1.0;
}`)
	vShader.Compile()
	fmt.Println("Vertex Shader")
	fmt.Println(vShader.Get(gl.COMPILE_STATUS))
	fmt.Println(vShader.GetInfoLog())

	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	defer fShader.Delete()
	fShader.Source(`#version 130
out vec3 color;

void main() {
  color = vec3(1,0,0);
}`)
	fShader.Compile()
	fmt.Println("Fragment Shader")
	fmt.Println(fShader.Get(gl.COMPILE_STATUS))
	fmt.Println(fShader.GetInfoLog())

	prog := gl.CreateProgram()
	prog.AttachShader(vShader)
	prog.AttachShader(fShader)
	prog.Link()
	fmt.Println("Program")
	fmt.Println(prog.Get(gl.LINK_STATUS))
	fmt.Println(prog.GetInfoLog())

	ta.prog = prog
}

func init() {
	Ctx = &glex.Context{
		Width:       800,
		Height:      600,
		ColorDepth:  [4]int{8, 8, 8, 8},
		BufferDepth: [2]int{0, 0},
		Vsync:       true,
		WindowTitle: "Experiment: Three",
		Keyboard: func(key, state int) {
			if key == glfw.KeyEsc {
				Ctx.Active = false
			}
		},
	}

	Window = &glex.GlfwContext{
		Ctx: Ctx,
	}

	Current = &TriActivity{}
	Current.SetContext(Ctx)
}

func main() {
	err := Window.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Window.Swap()
	gl.GenVertexArray().Bind()

	for Ctx.Active {
		Current.Draw()
		Window.Swap()
	}

}
