package main

import (
	"fmt"
	"os"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

var (
	running = true
	mouse   [3]int
)

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	err = glfw.OpenWindow(
		640,     // Width
		480,     // Height
		8, 8, 8, // RGB Bit depth
		8,             // Alpha Bit Depth
		0,             // Depth Buffer
		0,             // Stencil Buffer
		glfw.Windowed, // Window mode
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer glfw.CloseWindow()

	glfw.SetWindowTitle("Experiment: One")
	glfw.SetSwapInterval(1)
	glfw.SetKeyCallback(onKey)
	glfw.SetMouseButtonCallback(onMouseBtn)
	glfw.SetWindowSizeCallback(onResize)

	for running && glfw.WindowParam(glfw.Opened) == 1 {
		glfw.SwapBuffers()
	}
}

func onResize(w, h int) {
	gl.DrawBuffer(gl.FRONT_AND_BACK)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, w, h)
	gl.Ortho(0, float64(w), float64(h), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func onMouseBtn(button, state int) {
	mouse[button] = state
	fmt.Println(button, state)
}

func onKey(key, state int) {
	switch key {
	case glfw.KeyEsc:
		running = state == 0
	}
}
