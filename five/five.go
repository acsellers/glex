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
layout(location = 1) in vec3 vertexColor;

out vec3 fragmentColor;
uniform mat4 MVP;

void main() {
  gl_Position = MVP * vec4(vertexPosition_modelspace, 1);
  fragmentColor = vertexColor;
}`
	fragment = `#version 130
in vec3 fragmentColor;

out vec3 color;

void main() {
  color = fragmentColor;
}`
)

var (
	vertexData = []float32{
		-1.0, -1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, -1.0,
		1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
		1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		-1.0, 1.0, -1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, 1.0, 1.0,
		1.0, 1.0, 1.0,
		-1.0, 1.0, 1.0,
		1.0, -1.0, 1.0,
	}
	colorData = []float32{
		0.583, 0.771, 0.014,
		0.609, 0.115, 0.436,
		0.327, 0.483, 0.844,
		0.822, 0.569, 0.201,
		0.435, 0.602, 0.223,
		0.310, 0.747, 0.185,
		0.597, 0.770, 0.761,
		0.559, 0.436, 0.730,
		0.359, 0.583, 0.152,
		0.483, 0.596, 0.789,
		0.559, 0.861, 0.639,
		0.195, 0.548, 0.859,
		0.014, 0.184, 0.576,
		0.771, 0.328, 0.970,
		0.406, 0.615, 0.116,
		0.676, 0.977, 0.133,
		0.971, 0.572, 0.833,
		0.140, 0.616, 0.489,
		0.997, 0.513, 0.064,
		0.945, 0.719, 0.592,
		0.543, 0.021, 0.978,
		0.279, 0.317, 0.505,
		0.167, 0.620, 0.077,
		0.347, 0.857, 0.137,
		0.055, 0.953, 0.042,
		0.714, 0.505, 0.345,
		0.783, 0.290, 0.734,
		0.722, 0.645, 0.174,
		0.302, 0.455, 0.848,
		0.225, 0.587, 0.040,
		0.517, 0.713, 0.338,
		0.053, 0.959, 0.120,
		0.393, 0.621, 0.362,
		0.673, 0.211, 0.457,
		0.820, 0.883, 0.371,
		0.982, 0.099, 0.879,
	}
)

var Window glex.ContextImplementation

type CubeActivity struct {
	ctx    *glex.Context
	vbuf   gl.Buffer
	cbuf   gl.Buffer
	mvp    *math3d.Matrix
	matrix gl.UniformLocation
}

func (ca *CubeActivity) SetContext(c *glex.Context) {
	ca.ctx = c
}

func (ca *CubeActivity) SetActivity(c *glex.Context) {
	ca.ctx = c
}

func (ca *CubeActivity) Draw() {
	ca.ctx.Program.Use()

	ca.matrix.UniformMatrix4fv(false, ca.mvp.Values())

	vl := gl.AttribLocation(0)
	vl.EnableArray()
	ca.vbuf.Bind(gl.ARRAY_BUFFER)
	vl.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil,
	)

	cl := gl.AttribLocation(1)
	cl.EnableArray()
	ca.cbuf.Bind(gl.ARRAY_BUFFER)
	cl.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		nil,
	)
	gl.DrawArrays(gl.TRIANGLES, 0, 12*3)
	vl.DisableArray()
	cl.DisableArray()
}

func (ca *CubeActivity) Camera() {
}

func (ca *CubeActivity) Generate() {
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.GenVertexArray().Bind()

	ca.cbuf = gl.GenBuffer()
	ca.cbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(colorData)*4,
		colorData,
		gl.STATIC_DRAW,
	)

	ca.vbuf = gl.GenBuffer()
	ca.vbuf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(vertexData)*4,
		vertexData,
		gl.STATIC_DRAW,
	)

	proj := math3d.Perspective(45.0, 4.0/3.0, 0.1, 100.0)
	view := math3d.LookAt(
		math3d.NewVector([3]float32{4.0, 3.0, 3.0}),
		math3d.NewVector([3]float32{0, 0, 0}),
		math3d.NewVector([3]float32{0, 1.0, 0}),
	)
	model := math3d.Identity()

	ca.mvp = proj.Multiply(view.Multiply(model))

	ca.matrix = ca.ctx.Program.GetUniformLocation("MVP")
}

func init() {
	Window = &glex.GlfwContext{
		Ctx: &glex.Context{
			Width:       800,
			Height:      600,
			ColorDepth:  [4]int{8, 8, 8, 8},
			BufferDepth: [2]int{0, 0},
			Vsync:       true,
			WindowTitle: "Experiment: Five",
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
		Act: &CubeActivity{},
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
