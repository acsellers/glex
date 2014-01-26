package main

import (
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
	ctx *glex.Context
	buf gl.Buffer
	gen bool
}

func (ta *TriActivity) SetContext(c *glex.Context) {
	ta.ctx = c
}

func (ta *TriActivity) Draw() {
	ta.BlankScreen()
	ta.DrawTriangle()
}

func (ta *TriActivity) BlankScreen() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, ta.ctx.Width, ta.ctx.Height)
	gl.Ortho(0, float64(ta.ctx.Width), float64(ta.ctx.Height), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
}

func (ta *TriActivity) DrawTriangle() {
	if !ta.gen {
		ta.GenBuffer()
	}

	al := gl.AttribLocation(0)
	al.EnableArray()
	ta.buf.Bind(gl.ARRAY_BUFFER)
	al.AttribPointer(
		3,
		gl.FLOAT,
		false,
		0,
		0,
	)
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
	Window.Start()
	Window.Swap()
	gl.GenVertexArray().Bind()

	for Ctx.Active {
		Current.Draw()
		Window.Swap()
	}

}
