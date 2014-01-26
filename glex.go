package glex

import (
	"fmt"

	"github.com/go-gl/gl"
	"github.com/go-gl/glfw"
)

type Context struct {
	Width, Height int
	ColorDepth    [4]int
	BufferDepth   [2]int
	Fullscreen    bool
	Vsync         bool
	WindowTitle   string
	Active        bool

	Keyboard     func(key, state int)
	MouseButton  func(button, state int)
	MouseMove    func(x, y int)
	MouseWheel   func(pos int)
	WindowClose  func() int
	WindowResize func(width, height int)
}

type ContextImplementation interface {
	Start() error
	Swap()
	Refresh() error
	Close()
}

type GlfwContext struct {
	Ctx *Context
}

func (gc *GlfwContext) Start() error {
	c := gc.Ctx
	if err := glfw.Init(); err != nil {
		return err
	}
	fmt.Println(gl.Init())

	mode := glfw.Windowed
	if c.Fullscreen {
		mode = glfw.Fullscreen
	}

	err := glfw.OpenWindow(
		c.Width,
		c.Height,
		c.ColorDepth[0],
		c.ColorDepth[1],
		c.ColorDepth[2],
		c.ColorDepth[3],
		c.BufferDepth[0],
		c.BufferDepth[1],
		mode,
	)
	if err != nil {
		glfw.Terminate()
		return err
	}

	swapInterval := 0
	if c.Vsync {
		swapInterval = 1
	}
	glfw.SetSwapInterval(swapInterval)

	glfw.SetWindowTitle(c.WindowTitle)

	if c.Keyboard != nil {
		glfw.SetKeyCallback(c.Keyboard)
	}

	if c.MouseButton != nil {
		glfw.SetMouseButtonCallback(c.MouseButton)
	}

	if c.MouseMove != nil {
		glfw.SetMousePosCallback(c.MouseMove)
	}

	if c.MouseWheel != nil {
		glfw.SetMouseWheelCallback(c.MouseWheel)
	}

	glfw.SetWindowCloseCallback(func() int {
		gc.Ctx.Active = false
		if c.WindowClose != nil {
			return c.WindowClose()
		}
		return 0
	})

	if c.WindowResize != nil {
		glfw.SetWindowSizeCallback(c.WindowResize)
	}

	c.Active = true
	return nil
}

func (c *GlfwContext) Swap() {
	glfw.SwapBuffers()
}
func (c *GlfwContext) Refresh() error {
	return nil
}

func (c *GlfwContext) Close() {
	glfw.CloseWindow()
	glfw.Terminate()
}

type Activity interface {
	Draw()
	SetContext(*Context)
}

type NullActivity struct {
	Ctx *Context
}

func (na *NullActivity) SetContext(c *Context) {
	na.Ctx = c
}

func (na *NullActivity) Draw() {
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Viewport(0, 0, na.Ctx.Width, na.Ctx.Height)
	gl.Ortho(0, float64(na.Ctx.Width), float64(na.Ctx.Height), 0, -1.0, 1.0)
	gl.ClearColor(1, 1, 1, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
}
