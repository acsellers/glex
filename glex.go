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
	Shaders       [2]string
	Program       gl.Program

	Keyboard     func(key, state int)
	MouseButton  func(button, state int)
	MouseMove    func(x, y int)
	MouseWheel   func(pos int)
	WindowClose  func() int
	WindowResize func(width, height int)
}

func (c *Context) CompileShaders() {
	vShader := gl.CreateShader(gl.VERTEX_SHADER)
	defer vShader.Delete()
	vShader.Source(c.Shaders[0])
	vShader.Compile()
	if vShader.Get(gl.COMPILE_STATUS) != 1 {
		fmt.Println("Vertex Shader")
		fmt.Println(vShader.GetInfoLog())
	}

	fShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	defer fShader.Delete()
	fShader.Source(c.Shaders[1])
	fShader.Compile()
	if fShader.Get(gl.COMPILE_STATUS) != 1 {
		fmt.Println("Fragment Shader")
		fmt.Println(fShader.GetInfoLog())
	}

	prog := gl.CreateProgram()
	prog.AttachShader(vShader)
	prog.AttachShader(fShader)
	prog.Link()
	if prog.Get(gl.LINK_STATUS) != 1 {
		fmt.Println("Program")
		fmt.Println(prog.GetInfoLog())
	}

	c.Program = prog
}

type ContextImplementation interface {
	Start() error
	Activity() Activity
	SetActivity(Activity)
	Swap()
	Active() bool
	Refresh() error
	Close()
	Context() *Context
}

type Activity interface {
	Generate()
	Draw()
	Camera()
	SetContext(*Context)
}

type GlfwContext struct {
	Ctx *Context
	Act Activity
}

func (gc *GlfwContext) Active() bool {
	return gc.Ctx.Active
}

func (gc *GlfwContext) Start() error {
	c := gc.Ctx
	if err := glfw.Init(); err != nil {
		return err
	}

	mode := glfw.Windowed
	if c.Fullscreen {
		mode = glfw.Fullscreen
	}

	glfw.OpenWindowHint(glfw.FsaaSamples, 4)
	glfw.OpenWindowHint(glfw.OpenGLVersionMajor, 3)
	glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 0)
	//glfw.OpenWindowHint(glfw.OpenGLVersionMinor, 3)
	//glfw.OpenWindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

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

	if e := gl.Init(); e != 0 {
		return fmt.Errorf("GL Init returned error code: %v", e)
	}
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

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

	if c.Shaders != [2]string{"", ""} {
		fmt.Println("compiling shaders")
		c.CompileShaders()
	}

	if gc.Act != nil {
		gc.Act.SetContext(c)
		gc.Act.Generate()
	}

	c.Active = true
	return nil
}

func (c *GlfwContext) Swap() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	c.Act.Draw()
	c.Act.Camera()

	glfw.SwapBuffers()
}
func (c *GlfwContext) Refresh() error {
	return nil
}

func (c *GlfwContext) Activity() Activity {
	return c.Act
}

func (c *GlfwContext) SetActivity(a Activity) {
	c.Act = a
	c.Act.SetContext(c.Ctx)
	c.Act.Generate()
}

func (c *GlfwContext) Close() {
	glfw.CloseWindow()
	glfw.Terminate()
}

func (c *GlfwContext) Context() *Context {
	return c.Ctx
}
