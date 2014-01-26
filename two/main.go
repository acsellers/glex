package main

import "github.com/go-gl/glfw"

var (
	Window  ContextImplementation
	Ctx     *Context
	Current Activity
)

func init() {
	Ctx = &Context{
		Width:       800,
		Height:      600,
		ColorDepth:  [4]int{8, 8, 8, 8},
		BufferDepth: [2]int{0, 0},
		Vsync:       true,
		WindowTitle: "Experiment: Two",
		Keyboard: func(key, state int) {
			if key == glfw.KeyEsc {
				Ctx.Active = false
			}
		},
	}

	Window = &GlfwContext{
		Ctx: Ctx,
	}

	Current = &NullActivity{}
}

func main() {
	Window.Start()
	for Ctx.Active {
		Current.Draw()
		Window.Swap()
	}
}
