package main

import "github.com/go-gl/gl"

type Activity interface {
	Draw()
}

type NullActivity struct{}

func (na *NullActivity) Draw() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, Ctx.Width, Ctx.Height)
	gl.Ortho(0, float64(Ctx.Width), float64(Ctx.Height), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	na.DrawLines()
}

func (na *NullActivity) DrawLines() {
	gl.Enable(gl.BLEND)
	gl.Enable(gl.POINT_SMOOTH)
	gl.Enable(gl.LINE_SMOOTH)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.Color4f(0.0, 0.0, 0.0, 0.1)
	gl.Begin(gl.LINES)
	gl.Vertex2i(10, 10)
	gl.Vertex2i(20, 10)

	gl.Vertex2i(30, 30)
	gl.Vertex2i(40, 40)

	gl.Vertex2i(90, 80)
	gl.Vertex2i(70, 20)

	gl.End()
}
